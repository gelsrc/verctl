//
// Copyright (c) 2016 ЗАО Геликон Про http://www.gelicon.biz
//
package version

import (
	"testing"
)

func testStart(t *testing.T, level int, exp string) {
	if s := (&Ver{}).Start(level).Render(); s != exp {
		t.Errorf("Start(%v), got [%v], expected [%v]", level, s, exp)
	}
}

func TestStart(t *testing.T) {
	testStart(t, 1, "1-SNAPSHOT")
	testStart(t, 2, "1.0-SNAPSHOT")
	testStart(t, 3, "1.0.0-SNAPSHOT")
}

func testParseRender(t *testing.T, inp string) {
	if s := (&Ver{}).Parse(inp).Render(); s != inp {
		t.Errorf("Parse or render error, got [%v], expected [%v]", s, inp)
	}
}

func TestParseRender(t *testing.T) {
	testParseRender(t, "")
	testParseRender(t, "1")
	testParseRender(t, "1.")
	testParseRender(t, ".1")
	testParseRender(t, "1.2")
	testParseRender(t, "1.2.3")
	testParseRender(t, "1.pre2.3")
	testParseRender(t, "pre.fix1.2.3")
	testParseRender(t, "1.2.3suf.fix")
}

func testNext(t *testing.T, inp string, exp string, level int) {
	if s := (&Ver{}).Parse(inp).Next(level).Render(); s != exp {
		t.Errorf("Next error for %v, source [%v], got [%v], expected [%v]", level, inp, s, exp)
	}
}

func TestNext(t *testing.T) {
	testNext(t, "", "", 0)
	testNext(t, "1", "1", 0)

	testNext(t, "", "1-SNAPSHOT", 1)
	testNext(t, "1", "2-SNAPSHOT", 1)
	testNext(t, "1-SNAPSHOT", "2-SNAPSHOT", 1)
	testNext(t, "1.0-SNAPSHOT", "2.0-SNAPSHOT", 1)
	testNext(t, "1.2-SNAPSHOT", "2.0-SNAPSHOT", 1)
	testNext(t, "1.2.3", "2.0.0-SNAPSHOT", 1)

	testNext(t, "", "0.1-SNAPSHOT", 2)
	testNext(t, "1.0-SNAPSHOT", "1.1-SNAPSHOT", 2)
	testNext(t, "1.2-SNAPSHOT", "1.3-SNAPSHOT", 2)
	testNext(t, "1.2", "1.3-SNAPSHOT", 2)
	testNext(t, "1.2.3", "1.3.0-SNAPSHOT", 2)

	testNext(t, "pre.fix.1.2.3", "0.1.pre.fix.1.2.3-SNAPSHOT", 2)
}

func testRelease(t *testing.T, inp string, exp string) {
	if s := (&Ver{}).Parse(inp).Release().Render(); s != exp {
		t.Errorf("Release error, source [%v], got [%v], expected [%v]", inp, s, exp)
	}
}

func TestRelease(t *testing.T) {
	testRelease(t, "", "")
	testRelease(t, "1.2-SNAPSHOT", "1.2")
	testRelease(t, "1.2", "1.3")
	testRelease(t, "1.2.3", "1.2.4")
}

func testSnapshot(t *testing.T, inp string, exp string) {
	if s := (&Ver{}).Parse(inp).Snapshot().Render(); s != exp {
		t.Errorf("Snapshot error, source [%v], got [%v], expected [%v]", inp, s, exp)
	}
}

func TestSnapshot(t *testing.T) {
	testSnapshot(t, "", "")
	testSnapshot(t, "1.0-SNAPSHOT", "1.1-SNAPSHOT")
	testSnapshot(t, "1.2-SNAPSHOT", "1.3-SNAPSHOT")
	testSnapshot(t, "1.2", "1.3-SNAPSHOT")
	testSnapshot(t, "1.2.3", "1.2.4-SNAPSHOT")
	testSnapshot(t, "1.2.3.4.5.6", "1.2.3.4.5.7-SNAPSHOT")
}

func testLevel(t *testing.T, inp string, exp int) {
	if i := (&Ver{}).Parse(inp).Level(); i != exp {
		t.Errorf("Snapshot error, source [%v], got [%v], expected [%v]", inp, i, exp)
	}
}

func TestLevel(t *testing.T) {
	testLevel(t, "", 0)
	testLevel(t, "pre.fix", 0)
	testLevel(t, "1", 1)
	testLevel(t, "1-SNAPSHOT", 1)
	testLevel(t, "1.2", 2)
	testLevel(t, "1.2.3suffix", 3)
	testLevel(t, "1.2.3.suf.fix", 3)
}

func testTrunk(t *testing.T, level int, inp string, exp string) {
	if s := (&Ver{}).Parse(inp).Trunk(level).Render(); s != exp {
		t.Errorf("Trunk error for %v, source [%v], got [%v], expected [%v]", level, inp, s, exp)
	}
}

func TestTrunk(t *testing.T) {
	testTrunk(t, 0, "1.2.3", "")
	testTrunk(t, 1, "1.2.3", "1")
	testTrunk(t, 1, "1a.2b.3c", "1a")
	testTrunk(t, 2, "1.2.3-SNAPSHOT", "1.2-SNAPSHOT")
	testTrunk(t, 3, "1.2", "1.2")
	testTrunk(t, 3, "a.b.c.d.e", "")
	testTrunk(t, 3, "1.2.a.b.c.d.e", "1.2")
}
