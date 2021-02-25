package switchbot

import (
	"context"
	"errors"
)

type SceneService struct {
	c *Client
}

func newSceneService(c *Client) *SceneService {
	return &SceneService{c: c}
}

func (c *Client) Scene() *SceneService {
	return c.sceneService
}

type scenesResponse struct {
	StatusCode int     `json:"statusCode"`
	Mesasge    string  `json:"message"`
	Body       []Scene `json:"body"`
}

type Scene struct {
	ID   string `json:"sceneId"`
	Name string `json:"sceneName"`
}

// List get a list of manual scenes created by the current user.
// The first returned value is a list of scenes.
func (svc *SceneService) List(ctx context.Context) ([]Scene, error) {
	const path = "/v1.0/scenes"

	resp, err := svc.c.get(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	var response scenesResponse
	if err := resp.DecodeJSON(&response); err != nil {
		return nil, err
	}

	if response.StatusCode == 190 {
		return nil, errors.New("device internal error due to device states not synchronized with server")
	}

	return response.Body, nil
}

type sceneExecuteResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// Execute sends a request to execute a manual scene.
// The first given argument `id` is a scene ID which you want to execute, which can
// be retrieved by (*Client).Scene().List() function.
func (svc *SceneService) Execute(ctx context.Context, id string) error {
	path := "/v1.0/scenes/" + id + "/execute"

	resp, err := svc.c.post(ctx, path, nil)
	if err != nil {
		return err
	}
	defer resp.Close()

	var response sceneExecuteResponse
	if err := resp.DecodeJSON(&response); err != nil {
		return err
	}

	if response.StatusCode == 190 {
		return errors.New("device internal error due to device states not synchronized with server")
	}

	return nil
}