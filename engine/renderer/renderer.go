package renderer

import (
	"sync"

	"github.com/spaghettifunk/alaska-engine/engine/core"
	"github.com/spaghettifunk/alaska-engine/engine/platform"
	"github.com/spaghettifunk/alaska-engine/engine/renderer/vulkan"
)

type RendererBackend interface {
	Initialize(appName string) error
	Shutdow() error
	Resized(width, height uint16) error
	BeginFrame(deltaTime float64) error
	EndFrame(deltaTime float64) error
}

type RendererType uint8

const (
	Vulkan RendererType = iota
	DirectX
	Metal
	OpenGL
)

type Renderer struct {
	backend RendererBackend
}

type RenderPacket struct {
	DeltaTime float64
}

var initRenderer sync.Once
var renderer *Renderer

func Initialize(appName string, platform *platform.Platform) error {
	initRenderer.Do(func() {
		renderer = &Renderer{
			backend: vulkan.New(platform),
		}
	})
	return renderer.backend.Initialize(appName)
}

func Shutdown() error {
	return renderer.backend.Shutdow()
}

func BeginFrame(deltaTime float64) error {
	return renderer.backend.BeginFrame(deltaTime)
}

func EndFrame(deltaTime float64) error {
	return renderer.backend.EndFrame(deltaTime)
}

func OnResize(width, height uint16) error {
	return renderer.backend.Resized(width, height)
}

func DrawFrame(renderPacket *RenderPacket) error {
	if err := BeginFrame(renderPacket.DeltaTime); err != nil {
		core.LogError(err.Error())
		return err
	}
	if err := EndFrame(renderPacket.DeltaTime); err != nil {
		core.LogError("RendererEndFrame failed. Application shutting down...")
		return err
	}
	return nil
}
