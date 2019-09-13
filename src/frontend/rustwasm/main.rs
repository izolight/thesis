use sha2::{Sha256, Digest};

fn prgoressiveHash(data: u32) {
   hasher.input(data);
}

fn startHash() {
    let mut hasher = Sha256::default();
}

fn getHash() {
    let hash = hasher.result();
    return hash;
}

fn main() {

}
