package main

import (
  "bufio"
  "fmt"
  "log"
  "os"
  "sort"
  "time"
  "github.com/fatih/color"
	ui "github.com/gizak/termui/v3"
  "github.com/gizak/termui/v3/widgets"
  "github.com/gosuri/uilive"
  "github.com/shomali11/util/xhashes"
)

type match struct {
  occ   int
  hash  string
  desc  string
  lines []int
}

var all []match

func display(writer *uilive.Writer, all []match) {
  sort.Slice(all, func(i, j int) bool {
		return all[i].occ > all[j].occ
	})
  var top int = 10
  topAll := all[:top]
  for _, a := range topAll {
    color.New(color.FgBlue).Fprintf(writer, "[%d]", a.occ)
    fmt.Fprintf(writer, "::")
    color.New(color.FgGreen).Fprintf(writer, "%s\n", a.desc)
  }
}

func data_list(m []match) []float64 {
  var list []float64
  for _, i := range m {
    list = append(list, float64(i.occ))
  }
  return list
} 

func label_list(m []match) []string {
  var list []string
  for _, i := range m {
    list = append(list, i.desc)
  }
  return list
} 

func table_result(all []match, lineMax int) {
  sort.Slice(all, func(i, j int) bool {
		return all[i].occ > all[j].occ
	})
  var top int = 10
  topAll := all[:top]

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

  data_buf := [10][200]float64{}
  size := lineMax / 200

  for k, t := range topAll {
    for _, l := range t.lines {
      if(l / size < 200) {
        data_buf[k][l / size]++
      }
    }
  }

	sl0 := widgets.NewSparkline()
	sl0.Data = data_buf[0][:]
	sl0.LineColor = ui.ColorGreen
	sl1 := widgets.NewSparkline()
	sl1.Data = data_buf[1][:]
	sl1.LineColor = ui.ColorRed
	sl2 := widgets.NewSparkline()
	sl2.Data = data_buf[2][:]
	sl2.LineColor = ui.ColorCyan
	sl3 := widgets.NewSparkline()
	sl3.Data = data_buf[3][:]
	sl3.LineColor = ui.ColorMagenta
	sl4 := widgets.NewSparkline()
	sl4.Data = data_buf[4][:]
	sl4.LineColor = ui.ColorYellow

	slg0 := widgets.NewSparklineGroup(sl0, sl1, sl2, sl3, sl4)
	slg0.Title = "Sparkline 0"
	slg0.SetRect(0, 0, 200, 40)

	ui.Render(slg0)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}

func main() {
  writer := uilive.New()
  writer.Start()
  reader := bufio.NewScanner(os.Stdin)
  over := make(chan bool)
  ticker := time.NewTicker(1 * time.Second)
  go func() {
    for {
      select {
      case <- over:
        return
      case <-ticker.C:
        display(writer, all)
      }
    }
  }()
  ln := 0

  for reader.Scan() {
    ln++
    line := reader.Text()
    h := xhashes.SHA1(line)
    if len(line) > 20 {
      line = line[0:20]
    }
    found := 0 
    for k, v := range all {
      if v.hash == h {
        all[k].occ++
        all[k].lines = append(all[k].lines, ln)
        found++
      }
    } 
    if found == 0 {
      var curr = match {
        occ: 1,
        hash: h,
        desc: line,
        lines: []int{ln},
      }
      all = append(all, curr)
    }
  }
  if reader.Err() != nil {
    os.Exit(1)
  }
  over <- true
  table_result(all, ln)
}
