package client

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/gptscript-ai/otto/pkg/api/types"
	v1 "github.com/gptscript-ai/otto/pkg/storage/apis/otto.gptscript.ai/v1"
)

func (c *Client) UpdateAgent(ctx context.Context, id string, manifest v1.AgentManifest) (*types.Agent, error) {
	_, resp, err := c.putJSON(ctx, fmt.Sprintf("/agents/%s", id), manifest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return toObject(resp, &types.Agent{})
}

func (c *Client) GetAgent(ctx context.Context, id string) (*types.Agent, error) {
	_, resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/agents/"+id), nil)
	if err != nil {
		return nil, err
	}

	return toObject(resp, &types.Agent{})
}

func (c *Client) CreateAgent(ctx context.Context, agent v1.AgentManifest) (*types.Agent, error) {
	_, resp, err := c.postJSON(ctx, fmt.Sprintf("/agents"), agent)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return toObject(resp, &types.Agent{})
}

type ListAgentsOptions struct {
	Slug string
}

func (c *Client) ListAgents(ctx context.Context, opts ...ListAgentsOptions) (result types.AgentList, err error) {
	defer func() {
		sort.Slice(result.Items, func(i, j int) bool {
			return result.Items[i].Metadata.Created.Before(result.Items[j].Metadata.Created)
		})
	}()

	var opt ListAgentsOptions
	for _, o := range opts {
		if o.Slug != "" {
			opt.Slug = o.Slug
		}
	}

	_, resp, err := c.doRequest(ctx, http.MethodGet, "/agents", nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	_, err = toObject(resp, &result)
	if err != nil {
		return result, err
	}

	if opt.Slug != "" {
		var filtered types.AgentList
		for _, agent := range result.Items {
			if agent.Slug == opt.Slug && agent.SlugAssigned {
				filtered.Items = append(filtered.Items, agent)
			}
		}
		result = filtered
	}

	return
}

func (c *Client) DeleteAgent(ctx context.Context, id string) error {
	_, resp, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/agents/"+id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}