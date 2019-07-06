package client

import javafx.application.Application
import javafx.stage.Stage

class MainWindow: Application(){
    override fun start(primaryStage: Stage?){
        if(primaryStage != null) {
            primaryStage.title = "Offline Signature Validation"
            primaryStage.show()
        }
    }

    companion object{
        @JvmStatic
        fun main(args: Array<String>){
            launch(MainWindow::class.java, *args)
        }
    }
}