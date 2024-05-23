package urlextract

import (
	"fmt"
	"io"
	"log"
	"strings"
	"testing"
)

func init() {
	log.SetOutput(io.Discard)
}

var table = []struct {
	input string
}{
	{input: "asd"},
	{input: "asd asd"},
	{input: "https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/\nhttps://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/\nhttps://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/\nhttps://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/\n"},
	{input: "asd https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 asd https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ bsd"},
	{input: "asd https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 asd https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ bsd asd https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 asd https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ bsd asd https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 asd https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ bsd asd https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 asd https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ bsd asd https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 asd https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ bsd asd https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 asd https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ bsd asd https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 asd https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ bsd asd https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 asd https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ bsd asd https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 asd https://www.tiktok.com/@flakeysalt/video/6987860998568889606?is_from_webapp=1 https://streetkitchen.hu/kids/hus-es-haletelek/sult-csirkemelles-szendvics/ bsd "},
}

func BenchmarkExtractUrlsFromText(b *testing.B) {
	for _, v := range table {
		b.Run(fmt.Sprintf("words_in_message_%d", len(strings.Fields(v.input))), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ExtractUrlsFromText(v.input)
			}
		})
	}
}
