package client

class MyApplication: Application(){
    override fun start(primaryStage: Stage?){
        //You code here
    }

    companion object{
        @JvmStatic
        fun main(args: Array<String>){
            launch(MyApplication::class.java, *args)
        }
    }
}