cmd := exec.Command("sleep", "5")
err := cmd.Start()
if err != nil {
    log.Fatal(err)
}
log.Printf("Waiting for command to finish...")
err = cmd.Wait()
log.Printf("Command finished with error: %v", err)

done := make(chan error)
go func() {
    done <- cmd.Wait()
}()
select {
    case <-time.After(3 * time.Second):
        if err := cmd.Process.Kill(); err != nil {
            log.Fatal("failed to kill: ", err)
        }
        <-done // allow goroutine to exit
        log.Println("process killed")
    case err := <-done:
        log.Printf("process done with error = %v", err)
}

