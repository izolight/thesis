export class Validate {
    private static readonly nullmsg = "Error: object was null";
    public static notNull<T>(obj: T): obj is Exclude<T, null> {
        if (obj == null) {
            console.trace(this.nullmsg);
            throw new ReferenceError(this.nullmsg);
        }
        return true;
    }

    public static notUndefined<T>(obj: T): obj is Exclude<T, undefined> {
        if (obj === undefined) {
            throw new ReferenceError(`Error: Object ${obj} was undefined`);
        }
        return true;
    }

    public static notNullNotUndefined<T>(obj: T): obj is NonNullable<T> {
        return Validate.notNull(obj) && Validate.notUndefined(obj);
    }
}
