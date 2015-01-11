package lwgo

import "testing"

func TestLwTx(t *testing.T) {
    // type LwTx struct {
    //     setup bool
    //     Pin int
    //     Repeats int
    //     Onval, Offval int
    //     Translate int
    //     Period int
    // }

    var zeroed LwTx
    defaults := LwTx{
        setup: false,
        Pin: 0,
        Repeats: 0,
        Onval: 0,
        Offval: 0,
        Translate: 0,
        Period: 0,
    }

    if zeroed.setup != defaults.setup {
        t.Error("LwTx.setup == ", zeroed.setup, ", wanted: ", defaults.setup)
    }
    if zeroed.Pin != defaults.Pin {
        t.Error("LwTx.Pin == ", zeroed.Pin, ", wanted: ", defaults.Pin)
    }
    if zeroed.Repeats != defaults.Repeats {
        t.Error("LwTx.Repeats == ", zeroed.Repeats, ", wanted: ", defaults.Repeats)
    }
    if zeroed.Onval != defaults.Onval {
        t.Error("LwTx.Onval == ", zeroed.Onval, ", wanted: ", defaults.Onval)
    }
    if zeroed.Offval != defaults.Offval {
        t.Error("LwTx.Offval == ", zeroed.Offval, ", wanted: ", defaults.Offval)
    }
    if zeroed.Translate != defaults.Translate {
        t.Error("LwTx.Translate == ", zeroed.Translate, ", wanted: ", defaults.Translate)
    }
    if zeroed.Period != defaults.Period {
        t.Error("LwTx.Period == ", zeroed.Period, ", wanted: ", defaults.Period)
    }
}

func TestNewLwTx(t *testing.T) {
    // type LwTx struct {
    //     setup bool
    //     Pin int
    //     Repeats int
    //     Onval, Offval int
    //     Translate int
    //     Period int
    // }

    initd := NewLwTx()
    defaults := LwTx{
        setup: false,
        Pin: 3,
        Repeats: 10,
        Onval: 1,
        Offval: 0,
        Translate: 1,
        Period: 140,
    }

    if initd.setup != defaults.setup {
        t.Error("NewLwTx.setup == ", initd.setup, ", wanted: ", defaults.setup)
    }
    if initd.Pin != defaults.Pin {
        t.Error("NewLwTx.Pin == ", initd.Pin, ", wanted: ", defaults.Pin)
    }
    if initd.Repeats != defaults.Repeats {
        t.Error("NewLwTx.Repeats == ", initd.Repeats, ", wanted: ", defaults.Repeats)
    }
    if initd.Onval != defaults.Onval {
        t.Error("NewLwTx.Onval == ", initd.Onval, ", wanted: ", defaults.Onval)
    }
    if initd.Offval != defaults.Offval {
        t.Error("NewLwTx.Offval == ", initd.Offval, ", wanted: ", defaults.Offval)
    }
    if initd.Translate != defaults.Translate {
        t.Error("NewLwTx.Translate == ", initd.Translate, ", wanted: ", defaults.Translate)
    }
    if initd.Period != defaults.Period {
        t.Error("NewLwTx.Period == ", initd.Period, ", wanted: ", defaults.Period)
    }
}

func TestSetup(t *testing.T) {
    err := wiringPiSetup()
    if err != nil {
        t.Error("wiringPiSetup() failed")
    }
}

func TestSetupPins(t *testing.T) {
    initd := NewLwTx()
    initd.SetupPins()
    if initd.setup == false {
        t.Error("LwTx.setup == ", initd.setup, ", wanted: ", true)
    }
}
