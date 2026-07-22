package designtool

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type Game struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type InputParams struct {
	Name string `json:"name" jsonschema:"description=the name of name"`
}

func GetGameURL(_ context.Context, params *InputParams) (string, error) {
	GameSet := []Game{
		{Name: "原神", URL: "https://ys.mihoyo.com/tool"},
		{Name: "鸣潮", URL: "https://mc.kurogames.com/tool"},
		{Name: "明日方舟", URL: "https://ak.hypergryph.com/tool"},
	}
	for _, game := range GameSet {
		if game.Name == params.Name {
			return game.URL, nil
		}
	}
	return "", nil
}

func CreateTool() tool.InvokableTool {
	getGameTool := utils.NewTool(&schema.ToolInfo{
		Name: "get_game_url",
		Desc: "get a game url by name",
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{
				"name": &schema.ParameterInfo{
					Type:     schema.String,
					Desc:     "game's name",
					Required: true,
				},
			}),
	}, GetGameURL)
	return getGameTool
}
