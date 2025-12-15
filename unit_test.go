package opera_test

import (
	"testing"

	"github.com/colafanta/go-opera"
	"github.com/stretchr/testify/assert"
)

func TestUnitUniqueness(t *testing.T) {
	is := assert.New(t)
	is.True(opera.U.IsUnit())

	var temptUnit opera.Unit
	is.NotEqual(temptUnit, opera.U)
}
