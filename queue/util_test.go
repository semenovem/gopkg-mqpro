package queue

import (
  "github.com/stretchr/testify/assert"
  "testing"
)

func TestUtil(t *testing.T) {
  t.Run("unionProps", func(t *testing.T) {
    m1 := map[string]interface{}{
      "one": "one",
    }
    m2 := map[string]interface{}{
      "two": "two",
    }
    unionProps(m1, m2)

    p1, _ := m1["one"]
    p2, _ := m1["two"]
    assert.Equal(t, "one", p1)
    assert.Equal(t, "two", p2)

    p1, _ = m2["one"]
    p2, _ = m2["two"]
    assert.Nil(t, p1)
    assert.Equal(t, "two", p2)
  })

  t.Run("tailFour", func(t *testing.T) {
    d := []int{
      0, 1, 2, 3, 4, 5, 6, 7, 8,
    }
    r := []int{
      0, 3, 2, 1, 0, 3, 2, 1, 0,
    }

    for i, n := range d {
      assert.Equal(t, r[i], tailFour(n))
    }

    for i, n := range d {
      assert.Equal(t, int32(r[i]), tailFour32(int32(n)))
    }
  })
}
