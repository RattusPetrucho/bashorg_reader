package main

import "testing"

func Test_read_qoute(t *testing.T) {
    str, _ := read_qoute(1)

    if str != "<Ares> ppdv, все юниксы очень дружелюбны.. они просто очень разборчивы в друзьях ;)" {
        t.Error("did not expect got", str)
    }
}