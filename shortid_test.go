// Copyright (c) 2016 Ventu.io, Oleg Sklyar, contributors
// The use of this source code is governed by a MIT style license found in the LICENSE file

package shortid_test

import (
	"github.com/ventu-io/go-shortid"
	"testing"
	"time"
)

func Test_onGetDefault_defaultInstance(t *testing.T) {
	sid := shortid.GetDefault()
	expected := "Shortid(worker=0, epoch=2016-01-01 00:00:00 +0000 UTC, abc=Abc{alphabet='gzmZM7VINvOFcpho01x-fYPs8Q_urjq6RkiWGn4SHDdK5t2TAJbaBLEyUwlX9C3e'))"
	if sid.String() != expected {
		t.Errorf("expected %v", expected)
	}
}

func Test_onSetDefault_replacesDefaultInstance(t *testing.T) {
	expected := "Shortid(worker=0, epoch=2016-01-01 00:00:00 +0000 UTC, abc=Abc{alphabet='gzmZM7VINvOFcpho01x-fYPs8Q_urjq6RkiWGn4SHDdK5t2TAJbaBLEyUwlX9C3e'))"
	sid := shortid.GetDefault()
	if sid.String() != expected {
		t.Errorf("expected %v", expected)
	}
	// different worker, different seed (thus different shuffling)
	shortid.SetDefault(shortid.MustNew(1, shortid.DEFAULT_ABC, 2))
	expected = "Shortid(worker=1, epoch=2016-01-01 00:00:00 +0000 UTC, abc=Abc{alphabet='ip8bKduCDxnMQy-JrVHAN5h1s396jBvmFZOL0Pg2WTqwIE7f4ackXzoUSYlGt_eR'))"
	sid = shortid.GetDefault()
	if sid.String() != expected {
		t.Errorf("expected %v", expected)
	}
}

func Test_onGenerate_success(t *testing.T) {
	time.Sleep(2 * time.Millisecond)
	if id, err := shortid.Generate(); err != nil {
		t.Error(err)
	} else if len(id) != 9 {
		t.Error("incorrect id length")
	}
	time.Sleep(2 * time.Millisecond)
	if id, err := shortid.Generate(); err != nil {
		t.Error(err)
	} else if len(id) != 9 {
		t.Error("incorrect id length")
	}
}

func Test_oMustGenerate_success(t *testing.T) {
	time.Sleep(2 * time.Millisecond)
	if id := shortid.MustGenerate(); len(id) != 9 {
		t.Error("incorrect id length")
	}
	time.Sleep(2 * time.Millisecond)
	if id := shortid.MustGenerate(); len(id) != 9 {
		t.Error("incorrect id length")
	}
}

func TestShortid_onNew_success(t *testing.T) {
	expected := "Shortid(worker=5, epoch=2016-01-01 00:00:00 +0000 UTC, abc=Abc{alphabet='gzmZM7VINvOFcpho01x-fYPs8Q_urjq6RkiWGn4SHDdK5t2TAJbaBLEyUwlX9C3e'))"
	if sid, err := shortid.New(5, shortid.DEFAULT_ABC, 1); err != nil {
		t.Error(err)
	} else if sid.String() != expected {
		t.Errorf("expected %v, found %v", expected, sid.String())
	}
	expected = "Shortid(worker=31, epoch=2016-01-01 00:00:00 +0000 UTC, abc=Abc{alphabet='gzmZM7VINvOFcpho01x-fYPs8Q_urjq6RkiWGn4SHDdK5t2TAJbaBLEyUwlX9C3e'))"
	if sid, err := shortid.New(31, shortid.DEFAULT_ABC, 1); err != nil {
		t.Error(err)
	} else if sid.String() != expected {
		t.Errorf("expected %v, found %v", expected, sid.String())
	}
}

func TestShortid_onNew_withDifferentSeed_successWithDifferentShuffling(t *testing.T) {
	expected := "Shortid(worker=2, epoch=2016-01-01 00:00:00 +0000 UTC, abc=Abc{alphabet='gzmZM7VINvOFcpho01x-fYPs8Q_urjq6RkiWGn4SHDdK5t2TAJbaBLEyUwlX9C3e'))"
	if sid, err := shortid.New(2, shortid.DEFAULT_ABC, 1); err != nil {
		t.Error(err)
	} else if sid.String() != expected {
		t.Errorf("expected %v, found %v", expected, sid.String())
	}
	expected = "Shortid(worker=2, epoch=2016-01-01 00:00:00 +0000 UTC, abc=Abc{alphabet='U8dEc3Hnuq_RfyDApaT1ZxQmYePBCNMkF4-KJSvhjw609I7GlbzsriOL52XVoWgt'))"
	if sid, err := shortid.New(2, shortid.DEFAULT_ABC, 345234); err != nil {
		t.Error(err)
	} else if sid.String() != expected {
		t.Errorf("expected %v, found %v", expected, sid.String())
	}
}

func TestShortid_onNew_withWorkerAbove31_error(t *testing.T) {
	if _, err := shortid.New(32, shortid.DEFAULT_ABC, 1); err == nil {
		t.Error("expected error")
	}
}

func TestShortid_onNew_withIncorrectAbs_error(t *testing.T) {
	if _, err := shortid.New(1, "aasefvowefvjaHEFV", 1); err == nil {
		t.Error("expected error")
	}
}

func TestShortid_onMustNew_success(t *testing.T) {
	expected := "Shortid(worker=2, epoch=2016-01-01 00:00:00 +0000 UTC, abc=Abc{alphabet='gzmZM7VINvOFcpho01x-fYPs8Q_urjq6RkiWGn4SHDdK5t2TAJbaBLEyUwlX9C3e'))"
	if sid := shortid.MustNew(2, shortid.DEFAULT_ABC, 1); sid.String() != expected {
		t.Errorf("expected %v, found %v", expected, sid.String())
	}
}

func TestShortid_onError_panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	shortid.MustNew(1, "aasefvowefvjaHEFV", 1)
}

func TestShortid_onGenerate_success(t *testing.T) {
	sid := shortid.MustNew(1, shortid.DEFAULT_ABC, 1)
	if id, err := sid.Generate(); err != nil {
		t.Error(err)
	} else if len(id) != 9 {
		t.Errorf("expected id of length 9, found %v", id)
	}
	time.Sleep(2 * time.Millisecond)
	if id, err := sid.Generate(); err != nil {
		t.Error(err)
	} else if len(id) != 9 {
		t.Errorf("expected id of length 9, found %v", id)
	}
}

func TestShortid_onMustGenerate_success(t *testing.T) {
	sid := shortid.MustNew(1, shortid.DEFAULT_ABC, 1)
	if id := sid.MustGenerate(); len(id) != 9 {
		t.Errorf("expected id of length 9, found %v", id)
	}
	time.Sleep(2 * time.Millisecond)
	if id := sid.MustGenerate(); len(id) != 9 {
		t.Errorf("expected id of length 9, found %v", id)
	}
}

// func TestShortid_onGenerateInternal -- covered in integration tests

func TestShortid_onAbc_success(t *testing.T) {
	sid := shortid.MustNew(1, shortid.DEFAULT_ABC, 1)
	expected := "Abc{alphabet='gzmZM7VINvOFcpho01x-fYPs8Q_urjq6RkiWGn4SHDdK5t2TAJbaBLEyUwlX9C3e')"
	if abc := sid.Abc(); abc.String() != expected {
		t.Errorf("expected %v, found %v", expected, abc.String())
	}
}

func TestShortid_onEpoch_success(t *testing.T) {
	sid := shortid.MustNew(1, shortid.DEFAULT_ABC, 1)
	expected := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	if sid.Epoch() != expected {
		t.Errorf("expected %v, found %v", expected, sid.Epoch())
	}
}

func TestShortid_onWorker_success(t *testing.T) {
	sid := shortid.MustNew(25, shortid.DEFAULT_ABC, 1)
	if sid.Worker() != 25 {
		t.Errorf("expected 25, found %v", sid.Worker())
	}
}

func TestAbc_onNewAbc_success(t *testing.T) {
	expected := "Abc{alphabet='gzmZM7VINvOFcpho01x-fYPs8Q_urjq6RkiWGn4SHDdK5t2TAJbaBLEyUwlX9C3e')"
	if abc, err := shortid.NewAbc(shortid.DEFAULT_ABC, 1); err != nil {
		t.Error(err)
	} else if abc.String() != expected {
		t.Errorf("expected %v, found %v", expected, abc.String())
	}
}

func TestAbc_onNewAbc_wrongAlphabetLength_error(t *testing.T) {
	if _, err := shortid.NewAbc("asgliaeprugb", 1); err == nil {
		t.Error("expected error")
	}
	if _, err := shortid.NewAbc("1234567890qwertzuiopüäsdfghjklöä$<yxcvbnm,.->YXCVBNM;:_ASDFGHJKLQWERTZ", 1); err == nil {
		t.Error("expected error")
	}
}

func TestAbc_onNewAbc_alphabetNonUnique_error(t *testing.T) {
	runes := []rune(shortid.DEFAULT_ABC)
	if _, err := shortid.NewAbc(string(runes), 1); err != nil {
		t.Error(err)
	}
	runes[5] = 'A'
	if _, err := shortid.NewAbc(string(runes), 1); err == nil {
		t.Error("error expected")
	}
}

func TestAbc_onMustNewAbc_success(t *testing.T) {
	expected := "Abc{alphabet='gzmZM7VINvOFcpho01x-fYPs8Q_urjq6RkiWGn4SHDdK5t2TAJbaBLEyUwlX9C3e')"
	if abc := shortid.MustNewAbc(shortid.DEFAULT_ABC, 1); abc.String() != expected {
		t.Errorf("expected %v, found %v", expected, abc.String())
	}
}

func TestAbc_onMustNewAbc_onError_panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("panic expected")
		}
	}()
	shortid.MustNewAbc("asgliaeprugb", 1)
}

func TestAbc_onEncode_withVal0_success(t *testing.T) {
	abc := shortid.MustNewAbc(shortid.DEFAULT_ABC, 1)
	if id, err := abc.Encode(0, 1, 4); err != nil {
		t.Error(err)
	} else if len(id) != 1 {
		t.Errorf("expected 1 symbol")
	}
}

func TestAbc_onEncode_withValSmall_success(t *testing.T) {
	abc := shortid.MustNewAbc(shortid.DEFAULT_ABC, 1)
	if _, err := abc.Encode(48, 1, 4); err == nil {
		t.Errorf("expected error")
	}
	if id, err := abc.Encode(48, 2, 4); err != nil {
		t.Error(err)
	} else if len(id) != 2 {
		t.Errorf("expected 2 symbols")
	}
}

func TestAbc_onEncode_withValHuge_success(t *testing.T) {
	abc := shortid.MustNewAbc(shortid.DEFAULT_ABC, 1)
	if _, err := abc.Encode(214235345234524356, 14, 4); err == nil {
		t.Error("expected error")
	}
	if id, err := abc.Encode(214235345234524356, 15, 4); err != nil {
		t.Error(err)
	} else if len(id) != 15 {
		t.Errorf("expected 15 symbols")
	}
}

func TestAbc_onEncode_withNSymbols0_success(t *testing.T) {
	abc := shortid.MustNewAbc(shortid.DEFAULT_ABC, 1)
	if id, err := abc.Encode(214235345234524356, 0, 4); err != nil {
		t.Error(err)
	} else if len(id) != 15 {
		t.Errorf("expected 15 symbols")
	}
}

func TestAbc_onEncode_withDigits6_success(t *testing.T) {
	abc := shortid.MustNewAbc(shortid.DEFAULT_ABC, 1)
	if id, err := abc.Encode(214235345234524356, 0, 6); err != nil {
		t.Error(err)
	} else if len(id) != 10 {
		t.Errorf("expected 10 symbols")
	}
}

func TestAbc_onEncode_withDigitsWrong_error(t *testing.T) {
	abc := shortid.MustNewAbc(shortid.DEFAULT_ABC, 1)
	if _, err := abc.Encode(25, 0, 3); err == nil {
		t.Error("expected error")
	}
	if _, err := abc.Encode(25, 0, 4); err != nil {
		t.Error(err)
	}
	if _, err := abc.Encode(25, 0, 5); err != nil {
		t.Error(err)
	}
	if _, err := abc.Encode(25, 0, 6); err != nil {
		t.Error(err)
	}
	if _, err := abc.Encode(25, 0, 7); err == nil {
		t.Error("expected error")
	}
}

func TestAbc_onMustEncode_success(t *testing.T) {
	abc := shortid.MustNewAbc(shortid.DEFAULT_ABC, 1)
	if id := abc.MustEncode(25, 0, 4); len(id) != 2 {
		t.Errorf("expected len=2: %v", id)
	}
}

func TestAbc_onMustEncode_onError_panics(t *testing.T) {
	abc := shortid.MustNewAbc(shortid.DEFAULT_ABC, 1)
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	abc.MustEncode(25, 0, 2)

}

func TestAbc_onAlphabet_success(t *testing.T) {
	abc := shortid.MustNewAbc(shortid.DEFAULT_ABC, 1)
	expected := "gzmZM7VINvOFcpho01x-fYPs8Q_urjq6RkiWGn4SHDdK5t2TAJbaBLEyUwlX9C3e"
	if abc.Alphabet() != expected {
		t.Errorf("expected %v, found %v", expected, abc.Alphabet())
	}
}
