package test

import (
	"AesFileUtil/util"
	"reflect"
	"testing"
)

func TestAesDecrypt(t *testing.T) {
	type args struct {
		key  []byte
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := util.AesDecrypt(tt.args.key, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("AesDecrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AesDecrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAesEncrypt(t *testing.T) {
	type args struct {
		key  []byte
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := util.AesEncrypt(tt.args.key, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("AesEncrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AesEncrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAesFileEncode_decode(t *testing.T) {
	type fields struct {
		PwdKey []byte
	}
	type args struct {
		sourceFile      string
		destinationFile string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"文件 Aes 解密",
			fields{
				PwdKey: []byte("1234567811111111"),
			},
			args{sourceFile: "../dd.jpeg", destinationFile: "../abc.jpeg"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := util.AesFileEncode{
				PwdKey: tt.fields.PwdKey,
			}
			if err := a.Decode(tt.args.sourceFile, tt.args.destinationFile); (err != nil) != tt.wantErr {
				t.Errorf("decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAesFileEncode_encode(t *testing.T) {
	type fields struct {
		PwdKey []byte
	}
	type args struct {
		sourceFile      string
		destinationFile string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"文件 Aes 加密",
			fields{
				PwdKey: []byte("1234567812345678"),
			},
			args{sourceFile: "../WechatIMG44.jpeg", destinationFile: "../dd.jpeg"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := util.AesFileEncode{
				PwdKey: tt.fields.PwdKey,
			}
			if err := a.Encode(tt.args.sourceFile, tt.args.destinationFile); (err != nil) != tt.wantErr {
				t.Errorf("encode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecryptByAes(t *testing.T) {
	type args struct {
		key  []byte
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := util.DecryptByAes(tt.args.key, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptByAes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecryptByAes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecryptFile(t *testing.T) {
	type args struct {
		key             []byte
		sourceFile      string
		destinationFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := util.DecryptFile(tt.args.key, tt.args.sourceFile, tt.args.destinationFile); (err != nil) != tt.wantErr {
				t.Errorf("DecryptFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncryptByAes(t *testing.T) {
	type args struct {
		key  []byte
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := util.EncryptByAes(tt.args.key, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptByAes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EncryptByAes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncryptFile(t *testing.T) {
	type args struct {
		key             []byte
		sourceFile      string
		destinationFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := util.EncryptFile(tt.args.key, tt.args.sourceFile, tt.args.destinationFile); (err != nil) != tt.wantErr {
				t.Errorf("EncryptFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
