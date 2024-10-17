package metadata

import (
	"github.com/spaghettifunk/anima/engine/math"
	"github.com/spaghettifunk/anima/engine/resources"
)

const (
	BUILTIN_SHADER_NAME_SKYBOX   string = "Shader.Builtin.Skybox"
	BUILTIN_SHADER_NAME_MATERIAL string = "Shader.Builtin.Material"
	BUILTIN_SHADER_NAME_UI       string = "Shader.Builtin.UI"
)

/** @brief Represents a render target, which is used for rendering to a texture or set of textures. */
type RenderTarget struct {
	/** @brief Indicates if this render target should be updated on window resize. */
	SyncToWindowSize bool
	/** @brief The number of attachments */
	AttachmentCount uint8
	/** @brief An array of Attachments (pointers to textures). */
	Attachments []*resources.Texture
	/** @brief The renderer API internal framebuffer object. */
	InternalFramebuffer interface{}
}

/**
 * @brief The types of clearing to be done on a renderpass.
 * Can be combined together for multiple clearing functions.
 */
type RenderpassClearFlag uint32

const (
	/** @brief No clearing shoudl be done. */
	RENDERPASS_CLEAR_NONE_FLAG RenderpassClearFlag = 0x0
	/** @brief Clear the colour buffer. */
	RENDERPASS_CLEAR_COLOUR_BUFFER_FLAG RenderpassClearFlag = 0x1
	/** @brief Clear the depth buffer. */
	RENDERPASS_CLEAR_DEPTH_BUFFER_FLAG RenderpassClearFlag = 0x2
	/** @brief Clear the stencil buffer. */
	RENDERPASS_CLEAR_STENCIL_BUFFER_FLAG RenderpassClearFlag = 0x3
)

type RenderPassConfig struct {
	/** @brief The Name of this renderpass. */
	Name string
	/** @brief The name of the previous renderpass. */
	PrevName string
	/** @brief The name of the next renderpass. */
	NextName string
	/** @brief The current render area of the renderpass. */
	RenderArea math.Vec4
	/** @brief The clear colour used for this renderpass. */
	ClearColour math.Vec4
	/** @brief The clear flags for this renderpass. */
	ClearFlags uint8
}

/**
 * @brief Represents a generic RenderPass.
 */
type RenderPass struct {
	/** @brief The id of the renderpass */
	ID uint16
	/** @brief The current render area of the renderpass. */
	RenderArea math.Vec4
	/** @brief The clear colour used for this renderpass. */
	ClearColour math.Vec4
	/** @brief The clear flags for this renderpass. */
	ClearFlags uint8
	/** @brief The number of render targets for this renderpass. */
	RenderTargetCount uint8
	/** @brief An array of render Targets used by this renderpass. */
	Targets *RenderTarget
	/** @brief Internal renderpass data */
	InternalData interface{}
}

type RenderBufferType int

const (
	/** @brief Buffer is use is unknown. Default, but usually invalid. */
	RENDERBUFFER_TYPE_UNKNOWN RenderBufferType = iota
	/** @brief Buffer is used for vertex data. */
	RENDERBUFFER_TYPE_VERTEX
	/** @brief Buffer is used for index data. */
	RENDERBUFFER_TYPE_INDEX
	/** @brief Buffer is used for uniform data. */
	RENDERBUFFER_TYPE_UNIFORM
	/** @brief Buffer is used for staging purposes (i.e. from host-visible to device-local memory) */
	RENDERBUFFER_TYPE_STAGING
	/** @brief Buffer is used for reading purposes (i.e copy to from device local, then read) */
	RENDERBUFFER_TYPE_READ
	/** @brief Buffer is used for data storage. */
	RENDERBUFFER_TYPE_STORAGE
)

type RenderBuffer struct {
	/** @brief The type of buffer, which typically determines its use. */
	RenderBufferType RenderBufferType
	/** @brief The total size of the buffer in bytes. */
	TotalSize uint64
	/** @brief The amount of memory required to store the freelist. 0 if not used. */
	// freelist_memory_requirement uint64
	/** @brief The buffer freelist, if used. */
	// buffer_freelist freelist
	/** @brief The freelist memory block, if needed. */
	// freelist_block interface{}
	/** @brief Contains internal data for the renderer-API-specific buffer. */
	internal_data interface{}
}

/**
 * @brief A structure which is generated by the application and sent once
 * to the renderer to render a given frame. Consists of any data required,
 * such as delta time and a collection of views to be rendered.
 */
type RenderPacket struct {
	DeltaTime float64
	/** The number of views to be rendered. */
	ViewCount uint16
	/** An array of Views to be rendered. */
	Views []RenderViewPacket
}

/** @brief Known render view types, which have logic associated with them. */
type RenderViewKnownType int

const (
	/** @brief A view which only renders objects with *no* transparency. */
	RENDERER_VIEW_KNOWN_TYPE_WORLD RenderViewKnownType = 0x01
	/** @brief A view which only renders ui objects. */
	RENDERER_VIEW_KNOWN_TYPE_UI RenderViewKnownType = 0x02
	/** @brief A view which only renders skybox objects. */
	RENDERER_VIEW_KNOWN_TYPE_SKYBOX RenderViewKnownType = 0x03
)

/** @brief Known view matrix sources. */
type RenderViewViewMatrixSource int

const (
	RENDER_VIEW_VIEW_MATRIX_SOURCE_SCENE_CAMERA RenderViewViewMatrixSource = 0x01
	RENDER_VIEW_VIEW_MATRIX_SOURCE_UI_CAMERA    RenderViewViewMatrixSource = 0x02
	RENDER_VIEW_VIEW_MATRIX_SOURCE_LIGHT_CAMERA RenderViewViewMatrixSource = 0x03
)

/** @brief Known projection matrix sources. */
type RenderViewProjectionMatrixSource int

const (
	RENDER_VIEW_PROJECTION_MATRIX_SOURCE_DEFAULT_PERSPECTIVE  RenderViewProjectionMatrixSource = 0x01
	RENDER_VIEW_PROJECTION_MATRIX_SOURCE_DEFAULT_ORTHOGRAPHIC RenderViewProjectionMatrixSource = 0x02
)

/** @brief configuration for a renderpass to be associated with a view */
type RenderViewPassConfig struct {
	Name string
}

/**
 * @brief The configuration of a render view.
 * Used as a serialization target.
 */
type RenderViewConfig struct {
	/** @brief The Name of the view. */
	Name string
	/**
	 * @brief The name of a custom shader to be used
	 * instead of the view's default. Must be 0 if
	 * not used.
	 */
	CustomShaderName string
	/** @brief The Width of the view. Set to 0 for 100% Width. */
	Width uint16
	/** @brief The Height of the view. Set to 0 for 100% Height. */
	Height uint16
	/** @brief The known type of the view. Used to associate with view logic. */
	RenderViewType RenderViewKnownType
	/** @brief The source of the view matrix. */
	ViewMatrixSource RenderViewViewMatrixSource
	/** @brief The source of the projection matrix. */
	ProjectionMatrixSource RenderViewProjectionMatrixSource
	/** @brief The number of renderpasses used in this view. */
	PassCount uint8
	/** @brief The configuration of renderpasses used in this view. */
	Passes []RenderViewPassConfig
}

/**
 * @brief A render view instance, responsible for the generation
 * of view packets based on internal logic and given config.
 */
type RenderView struct {
	/** @brief The unique identifier of this view. */
	ID uint16
	/** @brief The Name of the view. */
	Name string
	/** @brief The current Width of this view. */
	Width uint16
	/** @brief The current Height of this view. */
	Height uint16
	/** @brief The known type of this view. */
	RenderViewType RenderViewKnownType
	/** @brief The number of renderpasses used by this view. */
	RenderpassCount uint8
	/** @brief An array of pointers to renderpasses used by this view. */
	Passes []*RenderPass
	/** @brief The name of the custom shader used by this view, if there is one. */
	CustomShaderName string
	/** @brief The internal, view-specific data for this view. */
	InternalData interface{}
	// callback functions
	OnCreate        OnCreateRenderView
	OnDestroy       OnDestroyRenderView
	OnResize        OnResizeRenderView
	OnBuildPacket   OnBuildPacketRenderView
	OnDestroyPacket OnDestroyPacketRenderView
	OnRender        OnRenderRenderView
}

/**
 * @brief A pointer to a function to be called when this view is created.
 *
 * @param self A pointer to the view being created.
 * @return True on success; otherwise false.
 */
type OnCreateRenderView func() bool

/**
 * @brief A pointer to a function to be called when this view is destroyed.
 *
 * @param self A pointer to the view being destroyed.
 */
type OnDestroyRenderView func() error

/**
 * @brief A pointer to a function to be called when the owner of this view (such
 * as the window) is resized.
 *
 * @param self A pointer to the view being resized.
 * @param width The new width in pixels.
 * @param width The new height in pixels.
 */
type OnResizeRenderView func(width, height uint32)

/**
 * @brief Builds a render view packet using the provided view and meshes.
 *
 * @param self A pointer to the view to use.
 * @param data Freeform data used to build the packet.
 * @param out_packet A pointer to hold the generated packet.
 * @return True on success; otherwise false.
 */
type OnBuildPacketRenderView func(data interface{}) (*RenderViewPacket, error)

/**
 * @brief Destroys the provided render view packet.
 * @param self A pointer to the view to use.
 * @param packet A pointer to the packet to be destroyed.
 */
type OnDestroyPacketRenderView func(packet *RenderViewPacket)

/**
 * @brief Uses the given view and packet to render the contents therein.
 *
 * @param self A pointer to the view to use.
 * @param packet A pointer to the packet whose data is to be rendered.
 * @param frame_number The current renderer frame number, typically used for data synchronization.
 * @param render_target_index The current render target index for renderers that use multiple render targets at once (i.e. Vulkan).
 * @return True on success; otherwise false.
 */
type OnRenderRenderView func(packet *RenderViewPacket, frame_number, render_target_index uint64) bool

/**
 * @brief A packet for and generated by a render view, which contains
 * data about what is to be rendered.
 */
type RenderViewPacket struct {
	/** @brief A constant pointer to the View this packet is associated with. */
	View *RenderView
	/** @brief The current view matrix. */
	ViewMatrix math.Mat4
	/** @brief The current projection matrix. */
	ProjectionMatrix math.Mat4
	/** @brief The current view position, if applicable. */
	ViewPosition math.Vec3
	/** @brief The current scene ambient colour, if applicable. */
	AmbientColour math.Vec4
	/** @brief The number of geometries to be drawn. */
	GeometryCount uint32
	/** @brief The Geometries to be drawn. */
	Geometries *GeometryRenderData
	/** @brief The name of the custom shader to use, if applicable. Otherwise 0. */
	CustomShadername string
	/** @brief Holds a pointer to freeform data, typically understood both by the object and consuming view. */
	ExtendedData interface{}
}

type GeometryRenderData struct {
	Model    math.Mat4
	Geometry *resources.Geometry
}

type MeshPacketData struct {
	MeshCount uint32
	Meshes    []*resources.Mesh
}

// type UIPacketData struct {
// 	MeshData *MeshPacketData
// 	// TODO: temp
// 	TextCount uint32
// 	Texts     []*resources.UIText
// }

type SkyboxPacketData struct {
	Skybox *resources.Skybox
}

/** @brief A range, typically of memory */
type MemoryRange struct {
	/** @brief The Offset in bytes. */
	Offset uint64
	/** @brief The size in bytes. */
	Size uint64
}
