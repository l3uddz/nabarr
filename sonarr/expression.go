package sonarr

import (
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/l3uddz/nabarr"
	"github.com/pkg/errors"
)

func (c *Client) compileExpressions(filters nabarr.PvrFilters) error {
	exprEnv := &nabarr.ExprEnv{}

	// compile ignores
	for _, ignoreExpr := range filters.Ignores {
		program, err := expr.Compile(ignoreExpr, expr.Env(exprEnv), expr.AsBool())
		if err != nil {
			return fmt.Errorf("ignore expression: %v: %w", ignoreExpr, err)
		}

		c.ignoresExpr = append(c.ignoresExpr, program)
	}

	return nil
}

func (c *Client) ShouldIgnore(mediaItem *nabarr.MediaItem) (bool, error) {
	exprItem := nabarr.GetExprEnv(mediaItem)

	for _, expression := range c.ignoresExpr {
		result, err := expr.Run(expression, exprItem)
		if err != nil {
			return true, fmt.Errorf("checking ignore expression: %w", err)
		}

		expResult, ok := result.(bool)
		if !ok {
			return true, errors.New("type assert ignore expression result")
		}

		if expResult {
			return true, nil
		}
	}

	return false, nil
}
