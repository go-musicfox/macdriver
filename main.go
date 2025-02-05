package main

import (
	"fmt"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/dispatch"
	"github.com/progrium/macdriver/objc"
	"os"
	"runtime"
	"time"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	go func() {
		Run()
		os.Exit(0)
	}()
	app := cocoa.NSApp()
	app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyProhibited)
	app.ActivateIgnoringOtherApps(true)
	app.Run()
}

func Run() {
	ok := make(chan bool)
	dispatch.Async(dispatch.MainQueue(), func() {
		//data, err := os.ReadFile("/Users/anhoder/Desktop/1.mp3")
		//fmt.Println(err)

		handlerCls := objc.NewClass("EventHandler", "NSObject")
		handlerCls.AddMethod("outputDeviceChanged:", func(notification objc.Object) {
			fmt.Println(notification)
		})
		objc.RegisterClass(handlerCls)

		handler := objc.Get("EventHandler").Alloc().Init()

		cls := objc.NewClass("SoundDelegate", "NSObject")
		cls.AddMethod("sound:didFinishPlaying:", func(sound objc.Object, didFinishPlaying bool) {
			fmt.Println("finish playing: ", didFinishPlaying)
			core.NSNotificationCenter_defaultCenter().
				AddObserver_selector_name_object_(handler, objc.Sel("outputDeviceChanged:"), core.String("test"), objc.Get("AVAudioSession").Get("sharedInstance"))
			if didFinishPlaying {
				ok <- true
			}
		})
		objc.RegisterClass(cls)

		delegate := objc.Get("SoundDelegate").Alloc().Init()

		url := core.NSURL_fileURLWithPath_isDirectory_(core.String("/Users/anhoder/Desktop/1.mp3"), false)
		s := cocoa.NSSound_InitWithURL(url, false)
		s.Set("delegate:", delegate)
		s.SetVolume_(0.8)
		fmt.Println(s.Volume())

		s.SetName_(core.String("test 111"))

		go func() {
			<-time.After(time.Second * 30)
			s.SetVolume_(1)
			fmt.Println(s.Volume())
		}()
		go func() {
			for {
				ticker := time.Tick(time.Second)
				<-ticker
				fmt.Println(s.Name(), s.CurrentTime(), s.Duration())
			}
		}()

		//d := core.NSData_WithBytes(data, uint64(len(data)))
		//s := cocoa.NSSound_InitWithData(d)
		s.Play()
	})
	<-ok

	//player := objc.Get("AVAudioPlayer").Alloc().Init()

	//player = player.Send("initWithContentsOfURL:", url)

}
