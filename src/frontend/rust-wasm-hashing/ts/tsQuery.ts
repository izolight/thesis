export function q(elementId: string): Exclude<HTMLElement, null> {
    const elem = document.getElementById(elementId);
    if (elem == null) {
        throw new ReferenceError(`Error: Element with id "${elementId}" does not exist`);
    } else {
        return elem;
    }
}
