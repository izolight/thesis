export class Validate {
    public static notNull<T>(obj: T): obj is Exclude<T, null> {
        if (obj == null) {
            throw new ReferenceError(`Error: Object ${obj} was null`);
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
