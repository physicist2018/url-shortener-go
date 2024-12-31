package urlstorage

import (
	"testing"
)

func TestURLStorage_CreateShortURL(t *testing.T) {
	type fields struct {
		Store []URLItem
	}
	type args struct {
		longURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				Store: []URLItem{
					{
						LongURL:  "yandex.ru",
						ShortURL: "wSv9wq",
					},
				},
			},
			args: args{
				longURL: "yandex.ru",
			},
			want:    "wSv9wq",
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				Store: GetDefaultURLStorage().Store,
			},
			args: args{
				longURL: "yandex.ru",
			},
			want:    "3TdGTj",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &URLStorage{
				Store: tt.fields.Store,
			}
			got, err := s.CreateShortURL(tt.args.longURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("URLStorage.CreateShortURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("URLStorage.CreateShortURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURLStorage_GetLongURL(t *testing.T) {
	type fields struct {
		Store []URLItem
	}
	type args struct {
		shortURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1 get existing key",
			fields: fields{
				Store: []URLItem{
					{
						LongURL:  "yandex.ru",
						ShortURL: "wSv9wq",
					},
				},
			},
			args: args{
				shortURL: "wSv9wq",
			},
			want:    "yandex.ru",
			wantErr: false,
		},
		{
			name: "test2 get not existing key",
			fields: fields{
				Store: []URLItem{
					{
						LongURL:  "yandex.ru",
						ShortURL: "wSv9w1",
					},
				},
			},
			args: args{
				shortURL: "wSv9wq",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &URLStorage{
				Store: tt.fields.Store,
			}
			got, err := s.GetLongURL(tt.args.shortURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("URLStorage.GetLongURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("URLStorage.GetLongURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
