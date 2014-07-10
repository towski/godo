package main

import(
    "github.com/tncardoso/gocurses"
    "math/rand"
    "time"
    "fmt"
    "runtime"
    "os"
    "github.com/towski/Golang-AStar/utils"
)

var disp int
var frames int
//var wg *sync.WaitGroup

func move(x int, y int) {
var temp int
    gocurses.Mvaddstr(y, x, " ")
    temp = rand.Intn(3)-1
    x += temp
    temp = rand.Intn(3)-1
    y += temp
    gocurses.Mvaddstr(y, x, "d")
    gocurses.Refresh()
    if(disp != 1){
 //       wg.Done()
    }
}

type Icon struct {
    x int
    y int
    old_x int
    old_y int
    char string
}

func main() {
    frames = 0
    messages := make(chan Icon)
    fmt.Print(time.Second)
    runtime.GOMAXPROCS(4)
    gocurses.Initscr()
    //defer gocurses.End()
    gocurses.Cbreak()
    gocurses.Noecho()
    gocurses.CursSet(0)
    gocurses.Stdscr.Keypad(true)
    var scene utils.Scene
    scene.InitScene(70, 100)
    scene.AddWalls(20)
    utils.InitAstar(&scene)


    for j := 0; j < 1; j++ {
        go func() { 
            var icon Icon
            var started = 0
            var finalPoint *utils.Point
            icon = Icon{}
            icon.x = 50//rand.Intn(50)
            icon.y = 50//rand.Intn(50)
            messages <- icon
            for {
                if started == 0 {
                    started = 1
                    scene.Data[3][3] = 'A'
                    scene.Data[5][5] = 'B'
                    for {
                        utils.FindPath(&scene)
                        //time.Sleep(50 * time.Millisecond)
                        if utils.Result != 10 {
                            break
                        }
                    }
                    icon.old_x = icon.x
                    icon.old_y = icon.y
                    finalPoint = &utils.FinalPoint
                    gocurses.End()
                    if finalPoint.Parent == nil {
                    fmt.Println("WHA")
                    } else {
                        icon.x = finalPoint.Parent.X
                        icon.y = finalPoint.Parent.Y
                        finalPoint = finalPoint.Parent
                        messages <- icon
                        time.Sleep(1000 * time.Millisecond)
                   }
                } else {
                    icon.old_x = icon.x
                    icon.old_y = icon.y
                        icon.x = finalPoint.Parent.X
                        icon.y = finalPoint.Parent.Y
                        finalPoint = finalPoint.Parent
                        messages <- icon
                        time.Sleep(1000 * time.Millisecond)
                }
            }
        }()
    }
                    //icon.old_x = x
                    //icon.old_y = y
                    //temp = rand.Intn(3)-1
                    //temp = rand.Intn(3)-1
                    //x += temp
                    //y += temp
                    //icon.x = x
                    //icon.y = y

    go func() { 
        for {
            icon := <-messages
            gocurses.Mvaddstr(icon.old_y, icon.old_x, " ")
            gocurses.Mvaddstr(icon.y, icon.x, "d")
            gocurses.Refresh()
            frames += 1
        }
    }()

    gocurses.Attron(gocurses.A_BOLD)
    disp = 1
    disp = gocurses.Getch()
     gocurses.End()
    fmt.Println(frames)
    os.Exit(0)
}
