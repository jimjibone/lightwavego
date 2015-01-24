package lwgo

import "testing"

func TestLwTx(t *testing.T) {
    // type LwTx struct {
    //     Pin int
    //     Repeats int
    //     Invert bool
    //     Translate bool
    //     Period int
    // }

    var zeroed LwTx
    defaults := LwTx{
        Pin: 0,
        Repeats: 0,
        Invert: false,
        Translate: false,
        Period: 0,
    }

    if zeroed.Pin != defaults.Pin {
        t.Error("LwTx.Pin == ", zeroed.Pin, ", wanted: ", defaults.Pin)
    }
    if zeroed.Repeats != defaults.Repeats {
        t.Error("LwTx.Repeats == ", zeroed.Repeats, ", wanted: ", defaults.Repeats)
    }
    if zeroed.Invert != defaults.Invert {
        t.Error("LwTx.Invert == ", zeroed.Invert, ", wanted: ", defaults.Invert)
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
    //     Pin int
    //     Repeats int
    //     Invert bool
    //     Translate bool
    //     Period int
    // }

    initd := NewLwTx()
    defaults := LwTx{
        Pin: 22,
        Repeats: 10,
        Invert: false,
        Translate: true,
        Period: 140,
    }

    if initd.Pin != defaults.Pin {
        t.Error("NewLwTx.Pin == ", initd.Pin, ", wanted: ", defaults.Pin)
    }
    if initd.Repeats != defaults.Repeats {
        t.Error("NewLwTx.Repeats == ", initd.Repeats, ", wanted: ", defaults.Repeats)
    }
    if initd.Invert != defaults.Invert {
        t.Error("NewLwTx.Invert == ", initd.Invert, ", wanted: ", defaults.Invert)
    }
    if initd.Translate != defaults.Translate {
        t.Error("NewLwTx.Translate == ", initd.Translate, ", wanted: ", defaults.Translate)
    }
    if initd.Period != defaults.Period {
        t.Error("NewLwTx.Period == ", initd.Period, ", wanted: ", defaults.Period)
    }
}

func TestNewLwBuffer( t *testing.T) {
    bytebuffer := make([]byte, 10)
    test, err := NewBuffer(bytebuffer)

    if err != nil {
        t.Error("NewBuffer returned with error: ", err)
    }

    if len(test.Raw()) != 10 {
        t.Error("NewBuffer length is too big. Length: ", len(test.Raw()))
    }

    for i, val := range test.Raw() {
        if bytebuffer[i] != val {
            t.Error("NewBuffer value at ", i, " is not equal to input: ", val, " != ", bytebuffer[i])
        }
    }
}
