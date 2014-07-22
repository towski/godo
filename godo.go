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
    "os/signal"
    "syscall"
)

var disp int
var frames int
//var wg *sync.WaitGroup

type Icon struct {
    x int
    y int
    old_x int
    old_y int
    char rune
}

func main() {
    frames = 0
    messages := make(chan Icon, 20)
    var log string
    fmt.Print(time.Second)
    runtime.GOMAXPROCS(4)
    gocurses.Initscr()
    defer gocurses.End()
    gocurses.Cbreak()
    gocurses.Noecho()
    gocurses.CursSet(0)
    gocurses.Stdscr.Keypad(true)
    var scene utils.Scene
    scene.InitScene(40, 40)
    scene.AddWalls(60)
    utils.InitAstar(&scene)
    var mutex = &sync.Mutex{}

    for j := 0; j < 1; j++ {
        go func() { 
            var icon Icon
            var target Icon
            var started = 0
            var finalPoint *utils.Point
            icon = Icon{}
            icon.x = utils.GetRandInt(20)
            icon.y = utils.GetRandInt(30)
            icon.char = 'd'
            messages <- icon
            time.Sleep(2000 * time.Millisecond)
            for {
                if started == 0 {
                    started = 1
                    target = Icon{}
                    target.x = utils.GetRandInt(20)
                    target.y = utils.GetRandInt(30)
                    target.char = 'f'
                    messages <- target
                    utils.SetOrig(&scene, icon.y, icon.x)
                    utils.SetDest(&scene, target.y, target.x)
                    for {
                        utils.FindPath(&scene)
                        //time.Sleep(50 * time.Millisecond)
                        if utils.Result != 10 {
                            break
                        }
                    }
                    if utils.Result == -1 {
                        started = 0
                        continue
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
                        scene.Data[icon.y][icon.x] = ' '
                        scene.Data[target.y][target.x] = ' '
                        utils.InitAstar(&scene)
                        time.Sleep(1000 * time.Millisecond)
                        started = 0
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
                gocurses.Mvaddch(draw_icon.old_y, draw_icon.old_x, ' ')
                log += fmt.Sprintf("removing d %d, %d\n", draw_icon.old_x, draw_icon.old_y)
            }
            gocurses.Mvaddch(draw_icon.y, draw_icon.x, draw_icon.char)
            log += fmt.Sprintf("Drawing d %d, %d %d\n", draw_icon.x, draw_icon.y, res)
            gocurses.Refresh()
            mutex.Unlock()
            frames += 1
        }
    }()

    go func() { 
        var wall_icon Icon
        wall_icon = Icon{}
        wall_icon.char = '#'
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

    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc,
        syscall.SIGHUP,
        syscall.SIGINT,
        syscall.SIGTERM,
        syscall.SIGQUIT)
    go func() {
        <-sigc
        gocurses.End()
        os.Exit(1)
        // ... do something ...
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
