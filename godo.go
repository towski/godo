package main

import(
    "github.com/tncardoso/gocurses"
//    "math/rand"
    "time"
    "fmt"
    "runtime"
    "os"
    "github.com/towski/Golang-AStar/utils"
    "sync"
)

var disp int
var frames int
//var wg *sync.WaitGroup

type Icon struct {
    x int
    y int
    old_x int
    old_y int
    char string
}

func main() {
    frames = 0
    messages := make(chan Icon, 20)
    var log string
    fmt.Print(time.Second)
    runtime.GOMAXPROCS(4)
    gocurses.Initscr()
    //defer gocurses.End()
    gocurses.Cbreak()
    gocurses.Noecho()
    gocurses.CursSet(0)
    gocurses.Stdscr.Keypad(true)
    var scene utils.Scene
    scene.InitScene(40, 40)
    scene.AddWalls(20)
    utils.InitAstar(&scene)
    var mutex = &sync.Mutex{}

    for j := 0; j < 1; j++ {
        go func() { 
            var icon Icon
            var started = 0
            var finalPoint *utils.Point
            icon = Icon{}
            icon.x = 20//rand.Intn(50)
            icon.y = 30//rand.Intn(50)
            messages <- icon
                        time.Sleep(2000 * time.Millisecond)
            for {
                if started == 0 {
                    started = 1
                    utils.SetOrig(&scene, 30, 20)
                    utils.SetDest(&scene, 3, 3)
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
                    for {
                        if finalPoint.Parent == nil {
                            break
                        } else {
                            finalPoint = finalPoint.Parent
                        }
                    }
                    gocurses.End()
                    if finalPoint.Child == nil {
                    } else {
                        icon.x = finalPoint.Child.Y
                        icon.y = finalPoint.Child.X
                        finalPoint = finalPoint.Child
                        messages <- icon
                        time.Sleep(300 * time.Millisecond)
                   }
                } else {
                    if finalPoint.Child != nil {
                        icon.old_x = icon.x
                        icon.old_y = icon.y
                        icon.x = finalPoint.Child.Y
                        icon.y = finalPoint.Child.X
                        finalPoint = finalPoint.Child
                        messages <- icon
                        time.Sleep(300 * time.Millisecond)
                    } else {
                        time.Sleep(1000 * time.Millisecond)
                    }
                }
            }
        }()
    }

    go func() { 
            var draw_icon Icon
            var res int
        for {
            draw_icon = <-messages
            mutex.Lock()
            if draw_icon.old_y != 0 && draw_icon.old_x != 0 {
                gocurses.Move(draw_icon.old_y, draw_icon.old_x)
                gocurses.Addch(' ')
                log += fmt.Sprintf("removing d %d, %d\n", draw_icon.old_x, draw_icon.old_y)
            }
            gocurses.Move(draw_icon.y, draw_icon.x)
            gocurses.Addch('d')
            log += fmt.Sprintf("Drawing d %d, %d %d\n", draw_icon.x, draw_icon.y, res)
            gocurses.Refresh()
            mutex.Unlock()
            frames += 1
        }
    }()

    go func() { 
        var wall_icon Icon
        wall_icon = Icon{}
        for i := 0; i < scene.Rows; i++ {
            for j := 0; j < scene.Cols; j++ {
                if scene.Data[i][j] == '#' {
             //       log += fmt.Sprintf("Drawing wall %d, %d\n", i, j)
                    wall_icon.x = j
                    wall_icon.y = i
                    messages <- wall_icon
                }
            }
        }
    }()

    gocurses.Attron(gocurses.A_BOLD)
    disp = 1
    disp = gocurses.Getch()
    gocurses.End()
    fmt.Println(log)
    fmt.Println(frames)
    scene.Draw()
    os.Exit(0)
}
