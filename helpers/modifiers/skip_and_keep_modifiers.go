package modifiers

import "github.com/vingarcia/kbuilder/ksqlmodifiers"

var skipInsertsModifier = ksqlmodifiers.AttrModifier{
	SkipOnInsert: true,
}

var skipUpdatesModifier = ksqlmodifiers.AttrModifier{
	SkipOnUpdate: true,
}

var nullableModifier = ksqlmodifiers.AttrModifier{
	Nullable: true,
}
