/* tslint:disable */
/**
*/
export class Sha256hasher {
  free(): void;
/**
* @returns {Sha256hasher} 
*/
  constructor();
/**
* @param {Uint8Array} input_bytes 
*/
  update(input_bytes: Uint8Array): void;
/**
* @returns {string} 
*/
  hex_digest(): string;
}
