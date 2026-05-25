package database

import (
	"testing"
)

func TestGetPokemon(t *testing.T) {
	Connect("mongodb+srv://admin:macacoverde@botdiscord.o7x8gjt.mongodb.net/")
	tests := []struct {
		// description of this test case
		// Named input parameters for target function.
		name string
		//want    *models.PokemonPool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Gardevoir",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := GetPokemon(tt.name)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetPokemon() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetPokemon() succeeded unexpectedly")
			}
			//// TODO: update the condition below to compare got with tt.want.
			//if true {
			//	t.Errorf("GetPokemon() = %v, want %v", got, tt.want)
			//}
		})
	}
}
