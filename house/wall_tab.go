package house

import (
  "github.com/runningwild/glop/gin"
  "github.com/runningwild/glop/gui"
  // "github.com/runningwild/glop/util/algorithm"
  "github.com/runningwild/haunts/base"
)

type WallPanel struct {
  *gui.VerticalTable
  room *roomDef
  viewer *RoomViewer

  wall_texture *WallTexture
  prev_wall_texture *WallTexture
  drag_anchor struct{ X,Y float32 }
  selected_walls map[int]bool
}

func MakeWallPanel(room *roomDef, viewer *RoomViewer) *WallPanel {
  var wp WallPanel
  wp.room = room
  wp.viewer = viewer
  wp.VerticalTable = gui.MakeVerticalTable()
  wp.selected_walls = make(map[int]bool)

  tex_table := gui.MakeVerticalTable()
  fnames := GetAllWallTextureNames()
  for i := range fnames {
    name := fnames[i]
    tex_table.AddChild(gui.MakeButton("standard", name, 300, 1, 1, 1, 1, func(t int64) {
      wt := MakeWallTexture(name)
      if wt == nil { return }
      wp.viewer.Temp.WallTexture = wt
      wp.viewer.Temp.WallTexture.X = 5
      wp.viewer.Temp.WallTexture.Y = 5
      wp.drag_anchor.X = 0
      wp.drag_anchor.Y = 0
    }))
  }
  wp.VerticalTable.AddChild(gui.MakeScrollFrame(tex_table, 300, 700))

  return &wp
}

func (w *WallPanel) textureNear(wx,wy int) *WallTexture {
  for _,tex := range w.room.WallTextures {
    var xx,yy float32
    if tex.X > float32(w.room.Size.Dx) {
      xx,yy,_ = w.viewer.modelviewToRightWall(float32(wx), float32(wy))
    } else if tex.Y > float32(w.room.Size.Dy) {
      xx,yy,_ = w.viewer.modelviewToLeftWall(float32(wx), float32(wy))
    } else {
      xx,yy,_ = w.viewer.modelviewToBoard(float32(wx), float32(wy))
    }
    dx := float32(tex.Texture.Data().Dx()) / 100 / 2
    dy := float32(tex.Texture.Data().Dy()) / 100 / 2
    if xx > tex.X - dx && xx < tex.X + dx && yy > tex.Y - dy && yy < tex.Y + dy {
      return tex
    }
  }
  return nil
}

func (w *WallPanel) Respond(ui *gui.Gui, group gui.EventGroup) bool {
  if w.VerticalTable.Respond(ui, group) {
    return true
  }
  if found,event := group.FindEvent(base.GetDefaultKeyMap()["flip"].Id()); found && event.Type == gin.Press {
    if w.wall_texture != nil {
      w.wall_texture.Flip = !w.wall_texture.Flip
    }
    return true
  }
  if found,event := group.FindEvent(gin.KeyDelete); found && event.Type == gin.Press {
    if w.wall_texture != nil {
      for i := range w.room.WallTextures {
        if w.room.WallTextures[i] == w.wall_texture {
          size := len(w.room.WallTextures)
          w.room.WallTextures[i] = w.room.WallTextures[size - 1]
          w.room.WallTextures = w.room.WallTextures[0 : size - 1]
          break
        }
      }
      w.wall_texture = nil
    }
    return true
  }
  if found,event := group.FindEvent(gin.MouseWheelVertical); found {
    if w.wall_texture != nil {
      w.wall_texture.Rot += float32(event.Key.CurPressAmt() / 100)
    }
  }
  if found,event := group.FindEvent(gin.MouseLButton); found && event.Type == gin.Press {
    if w.wall_texture != nil {
      w.wall_texture.temporary = false
      w.wall_texture = nil
    } else if w.wall_texture == nil {
      w.wall_texture = w.textureNear(event.Key.Cursor().Point())
      if w.wall_texture != nil {
        w.prev_wall_texture = new(WallTexture)
        *w.prev_wall_texture = *w.wall_texture
        w.wall_texture.temporary = true

        wx,wy := w.viewer.BoardToWindow(w.wall_texture.X, w.wall_texture.Y)
        px,py := event.Key.Cursor().Point()
        w.drag_anchor.X = float32(px) - float32(wx) - 0.5
        w.drag_anchor.Y = float32(py) - float32(wy) - 0.5
      }
    }
    return true
  }
  return false
}

func (w *WallPanel) Think(ui *gui.Gui, t int64) {
  if w.wall_texture != nil {
    px,py := gin.In().GetCursor("Mouse").Point()
    tx := float32(px) - w.drag_anchor.X
    ty := float32(py) - w.drag_anchor.Y
    bx,by := w.viewer.WindowToBoard(int(tx), int(ty))
    w.wall_texture.X = bx
    w.wall_texture.Y = by
  }
  w.VerticalTable.Think(ui, t)
}

func (w *WallPanel) Collapse() {
  if w.viewer.Temp.WallTexture != nil && w.prev_wall_texture != nil {
    w.room.WallTextures = append(w.room.WallTextures, w.prev_wall_texture)
  }
  w.prev_wall_texture = nil
  w.viewer.Temp.WallTexture = nil
}

func (w *WallPanel) Expand() {
}

func (w *WallPanel) Reload() {
}

