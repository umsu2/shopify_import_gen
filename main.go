package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// handle = prefix_adjective_noun_suffix 
// output file name
const le = `malesuada nunc vel risus commodo viverra maecenas accumsan lacus vel facilisis volutpat est velit egestas dui id ornare arcu odio ut sem nulla pharetra diam sit amet nisl suscipit adipiscing bibendum est ultricies integer quis auctor elit sed vulputate mi sit amet mauris commodo quis imperdiet massa tincidunt nunc pulvinar sapien et ligula ullamcorper malesuada proin libero nunc consequat interdum varius sit amet mattis vulputate enim nulla aliquet porttitor lacus luctus accumsan tortor posuere ac ut consequat semper viverra nam libero justo laoreet sit amet cursus sit amet dictum sit amet justo donec enim diam vulputate ut pharetra sit amet aliquam id diam maecenas ultricies mi eget mauris pharetra et ultrices neque ornare aenean euismod elementum nisi quis eleifend quam adipiscing vitae proin sagittis nisl rhoncus mattis rhoncus urna neque viverra justo nec ultrices dui sapien eget mi proin sed libero enim sed faucibus turpis in eu mi bibendum neque egestas congue quisque egestas diam in arcu cursus euismod quis viverra nibh cras pulvinar mattis nunc sed blandit libero volutpat sed cras ornare arcu dui vivamus arcu felis bibendum ut tristique et egestas quis ipsum suspendisse ultrices gravida dictum fusce ut placerat orci nulla pellentesque dignissim enim sit amet venenatis urna cursus eget nunc scelerisque viverra mauris in aliquam sem fringilla ut morbi tincidunt augue interdum velit euismod in pellentesque massa placerat duis ultricies lacus sed turpis tincidunt id aliquet risus feugiat in ante metus dictum at tempor commodo ullamcorper a lacus vestibulum sed arcu non odio euismod lacinia at quis risus sed vulputate odio ut enim blandit volutpat maecenas volutpat blandit aliquam etiam erat velit scelerisque in dictum non consectetur a erat nam at lectus urna duis convallis convallis tellus id interdum velit laoreet id donec ultrices tincidunt arcu non sodales neque sodales ut etiam sit amet nisl purus in mollis nunc sed id semper risus in hendrerit gravida rutrum quisque non tellus orci ac auctor augue mauris augue neque gravida in fermentum et sollicitudin ac orci phasellus egestas tellus rutrum tellus pellentesque eu tincidunt tortor aliquam nulla facilisi cras fermentum odio eu feugiat pretium nibh ipsum consequat nisl vel pretium lectus quam id leo in vitae turpis massa sed elementum tempus egestas sed sed risus pretium quam vulputate dignissim suspendisse in est ante in nibh mauris cursus mattis molestie a iaculis at erat pellentesque adipiscing commodo elit at imperdiet dui accumsan sit amet nulla facilisi morbi tempus iaculis urna id volutpat lacus laoreet non curabitur gravida arcu ac tortor dignissim convallis aenean et tortor at risus viverra adipiscing at in tellus integer feugiat scelerisque varius morbi enim nunc faucibus a pellentesque sit amet porttitor eget dolor morbi non arcu risus quis varius quam quisque id diam vel quam elementum pulvinar etiam non quam lacus suspendisse faucibus interdum posuere lorem ipsum dolor sit amet consectetur adipiscing elit duis tristique sollicitudin nibh sit amet commodo nulla facilisi nullam vehicula ipsum a arcu cursus vitae congue mauris rhoncus aenean vel elit scelerisque mauris pellentesque pulvinar pellentesque habitant morbi tristique senectus et`
const headers = "Handle,Title,Body (HTML),Vendor,Type,Tags,Published,Option1 Name,Option1 Value,Option2 Name,Option2 Value,Option3 Name,Option3 Value,Variant SKU,Variant Grams,Variant Inventory Tracker,Variant Inventory Qty,Variant Inventory Policy,Variant Fulfillment Service,Variant Price,Variant Compare At Price,Variant Requires Shipping,Variant Taxable,Variant Barcode,Image Src,Image Position,Image Alt Text,Gift Card,SEO Title,SEO Description,Google Shopping / Google Product Category,Google Shopping / Gender,Google Shopping / Age Group,Google Shopping / MPN,Google Shopping / AdWords Grouping,Google Shopping / AdWords Labels,Google Shopping / Condition,Google Shopping / Custom Product,Google Shopping / Custom Label 0,Google Shopping / Custom Label 1,Google Shopping / Custom Label 2,Google Shopping / Custom Label 3,Google Shopping / Custom Label 4,Variant Image,Variant Weight Unit,Variant Tax Code,Cost per item"
const maxBodyWords = 100
const (
	Handle = iota
	Title
	Body
	Vendor
	Type
	Tags
	Published
	Opt1Name
	Opt1Val
	InventoryPolicy              = 17
	InventoryFulfillmentService = 18
	VariantPrice                 = 19
	VariantCompareAtPrice        = 20
	CostPerItem                  = 46
)

const productPrefix = `generated_product`

var headerLength int
var rd *rand.Rand
var leWords []string

func main() {
	numOfProductsToGenerate := flag.Int("p", 1000, "# of products to generate")
	flag.Parse()

	if numOfProductsToGenerate == nil || *numOfProductsToGenerate == 0 {
		fmt.Println("missing required -p param")
		os.Exit(1)
	}

	writer := csv.NewWriter(os.Stdout)
	shopifyCSVHeaders := strings.Split(headers, ",")
	rd = rand.New(rand.NewSource(time.Now().Unix()))
	leWords = strings.Split(le, " ")
	headerLength = len(shopifyCSVHeaders)
	err := writer.Write(shopifyCSVHeaders)
	if err != nil {
		panic(err)
	}
	for i := 0; i < *numOfProductsToGenerate; i++ {
		words := []string{productPrefix, pickAdj(), pickNoun()}
		err = writer.Write(generateProdOptions(words, pickAdj(), time.Now()))
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()
}

func pickOne(in []string) string {
	if len(in) == 0 {
		panic("cannot pick from empty array")
	}
	idx := rd.Intn(len(in))
	return in[idx]
}

func generateProdOptions(itemDescriptions []string, option string, t time.Time) []string {
	handle := strings.Join(itemDescriptions, "_")
	handle = strings.Join([]string{handle, fmt.Sprintf("%d", t.Unix())}, "_")
	val := make([]string, headerLength)
	val[Handle] = handle
	val[Title] = strings.Join(itemDescriptions, " ")
	bs := rd.Intn(maxBodyWords)
	shuffle(leWords, bs)
	body := leWords[:bs]
	val[Body] = strings.Join(body, " ")
	val[Vendor] = "factory"
	val[Type] = "gizmos"
	val[Tags] = "fake,mock,hidden,magic_beans,do_not_use,generated"
	val[Published] = "True"
	val[InventoryPolicy] = "deny"
	val[InventoryFulfillmentService] = "manual"
	val[Opt1Name] = "Options"
	val[Opt1Val] = option
	val[VariantPrice] = fmt.Sprintf("%d.%d", rand.Intn(10000), rand.Intn(100))
	return val
}

func shuffle(vals []string, max int) {
	n := len(vals)
	for i := range vals {
		randIndex := rd.Intn(n)
		vals[i], vals[randIndex] = vals[randIndex], vals[i]
		if i >= max {
			break
		}
	}
}
