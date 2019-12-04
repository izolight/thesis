import {Callable, Continuable} from "./interfaces";

export class Queue {
    private readonly queue = Array<Callable>();
    private continuation: Callable;
    private started: Boolean = false;

    constructor(continuation: Callable) {
        this.continuation = continuation;
    }

    public add(c: Continuable) {
        if (this.started) {
            console.log("nope");
            return;
        }

        this.queue.push(() => {
            const next = this.queue.pop();
            if (next != undefined) {
                c(next);
            } else {
                c(this.continuation);
            }
        })
    }

    public start() {
        this.started = true;
        this.queue.reverse();
        const c = this.queue.pop();
        if (c != undefined) {
            c();
        }
    }
}
