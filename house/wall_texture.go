package house

import (
  "haunts/base"
  "haunts/texture"
  "github.com/arbaal/mathgl"
  "gl"
)

func init() {
  base.RegisterRegistry("wall textures", make(map[string]*wallTextureDef))
}

func MakeWallTexture(name string) *WallTexture {
  wt := WallTexture{ Defname: name }
  wt.Load()
  return &wt
}

func GetAllWallTextureNames() []string {
  return base.GetAllNamesInRegistry("wall textures")
}

func LoadAllWallTexturesInDir(dir string) {
  base.RegisterAllObjectsInDir("wall textures", dir, ".json", "json")
}

func (wt *WallTexture) Load() {
  base.LoadObject("wall textures", wt)
  if wt.wallTextureDef.texture_data == nil {
    wt.wallTextureDef.texture_data = texture.LoadFromPath(wt.wallTextureDef.Texture_path)
  }
}

type WallTexture struct {
  Defname string
  *wallTextureDef
  WallTextureInst
}

type wallTextureDef struct {
  // Name of this texture as it appears in the editor, should be unique among
  // all WallTextures
  Name string

  // Path to the texture - stored as a relative path but converted to an
  // absolute path when loaded
  Texture_path string `registry:"path"`

  // The texture itself
  texture_data *texture.Data
}

type WallTextureInst struct {
  // Position of the texture in floor coordinates.  If these coordinates exceed
  // either the dx or dy of the room, then this texture will be drawn, at least
  // partially, on the wall.  The coordinates should not both exceed the
  // dimensions of the room.
  X,Y float32
  Rot float32
}

func (wt *WallTexture) Render() {
  dx2 := float32(wt.texture_data.Dx) / 100 / 2
  dy2 := float32(wt.texture_data.Dy) / 100 / 2
  gl.Enable(gl.TEXTURE_2D)
  wt.texture_data.Bind()

  var rot mathgl.Mat3
  rot.RotationZ(wt.Rot)

  ll := mathgl.Vec2{ - dx2, - dy2 }
  ul := mathgl.Vec2{ - dx2, + dy2 }
  ur := mathgl.Vec2{ + dx2, + dy2 }
  lr := mathgl.Vec2{ + dx2, - dy2 }

  ll.Transform(&rot)
  ul.Transform(&rot)
  ur.Transform(&rot)
  lr.Transform(&rot)

  gl.Begin(gl.QUADS)
  gl.TexCoord2i(0, 0)
  gl.Vertex2f(wt.X + ll.X, wt.Y + ll.Y)
  gl.TexCoord2i(0, -1)
  gl.Vertex2f(wt.X + ul.X, wt.Y + ul.Y)
  gl.TexCoord2i(-1, -1)
  gl.Vertex2f(wt.X + ur.X, wt.Y + ur.Y)
  gl.TexCoord2i(-1, 0)
  gl.Vertex2f(wt.X + lr.X, wt.Y + lr.Y)
  gl.End()
}