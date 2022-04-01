package sonarr

import (
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/l3uddz/nabarr"
	"github.com/l3uddz/nabarr/media"
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

		c.ignoresExpr = append(c.ignoresExpr, nabarr.NewExprProgram(ignoreExpr, program))
	}

	return nil
}

func (c *Client) ShouldIgnore(mediaItem *media.Item) (bool, string, error) {
	exprItem := nabarr.NewExprEnv(mediaItem)

	for _, expression := range c.ignoresExpr {
		result, err := expr.Run(expression.Program, exprItem)
		if err != nil {
			return true, expression.String(), fmt.Errorf("checking ignore expression: %w", err)
		}

		expResult, ok := result.(bool)
		if !ok {
			return true, expression.String(), errors.New("type assert ignore expression result")
		}

		if expResult {
			return true, expression.String(), nil
		}
	}

	return false, "", nil
}
