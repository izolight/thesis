## Arbeitsjournal

#### 16.09.19
- Diskutieren der zu verwendenden Sprachen, Frameworks, Technologien
- Ausprobieren in-browser hashing
- Entwickeln Go+WASM in-browser hashing
- Anpassen BFH LaTeX template

#### 20.09.19
- Entwickeln Rust+WASM in-browser hashing
- Beginn Dokumentation
- Erstes Gespräch mit Betreuenden

#### 27.09.19
- Pflichtenheft erstellen

#### 30.09.19
- Weiterarbeiten am Pflichtenheft
- Zweites Gespräch mit Betreuenden

#### 04.10.19
- Weiterarbeiten am Pflichtenheft und der Dokumentation

#### 11.10.19
- Weiterarbeiten an der Dokumentation
- Evaluieren YubiHSM 2 und dokumentation

#### 12.10.19
- Weiterarbeiten an der Dokumentation
- Einleitung überarbeitet
- Weitere Requirements dokumentiert
- Einführung in die Kryptographie geschrieben

#### 14.10.19
- Drittes Gespräch mit Betreuenden

#### 18.10.19
- Weiterarbeiten an der Dokumentation
- Diverse UML Sequenzdiagramme erstellen
- Einführung in Krypto fertigstellen (Signaturen, PKIs, MITM)
- CSC Standard gelesen und Vergleich dokumentiert
- REST API Spezifikation erstellt

#### 23.10.19
- Lernen und Ausprobieren Ktor

#### 24.10.19
- Implementieren Beispiel-REST-Endpoint in Ktor

#### 27.10.19
- Implementieren Beispiel-REST-Endpoint in Golang

#### 28.10.19
- Rust WASM Hashing neu implementieren mit Rust 1.38.0

#### 29.10.19
- Vortrag für Experten vorbereiten
- Besprechung mit Betreuenden und Experten

#### 01.11.19
- Frontend UI layout sketches
- WASM Performance Vergleich fertigstellen
- UML Sequenzdiagramme verbessern
- Dokumentation
- Signaturdateiformat für LTV erweitern

#### 02.11.19
- Arbeiten am Verifier-Programm


#### 08.11.19
- Arbeiten an Backend-Implementation
- Signaturdateiformat für LTV erweitern
- Prozess der Entwicklung des Protokolls dokumentieren
- Dokumentation verbessern


#### 14.11.19
- Besprechung mit Betreuenden

#### 15.11.19
- Arbeiten am Verifier-Programm
    - beginn timestamp verifizierung
    - beginn ocsp verifizierung
    - beginn ltv verifierung
    - scripts um ocsp/timestamps zu generieren
- Arbeiten am Signing Server


#### 18.11.19
- Setup gitlab ci
    - template kopiert um pdf zu bauen(funktionier noch nicht)
- Beginn der Arbeit am Verifier-Programm
    - beginn id token verifizerung (test JWTs via okta)


#### 19.11.19
- Arbeit am Verifier-Programm
    - erste schritte um asynchron zu verifizieren
    - hash überprüfung

#### 20.11.19
- Arbeit am Verifier-Programm
    - optimierungen timestamp verifizierung

#### 22.11.19
- Arbeit am Signing Server: JWTs validieren
- Problem mit JWKS erkannt (kein X.509)
- Setup eigene CA mit OCSP Responder
    - pki mit rest api via cfssl
    - root & intermediate ca + ocsp responder
- Setup eigener OIDC IDP
    - eigener idp mit hydra

#### 27.11.19
- Arbeiten am Signing Server: 
    - Zertifikat und CSR generieren mit BouncyCastle
    - Anpassen OIDC Code an eigenen IdP
    - OIDC Tests mit eigenem IdP
- IdP fixen
    - hydra nur mässig brauchbar, da nur oidc provider ohne user verwaltung
    - rückbau von hydra und ersetzen mit keycloak

#### 28.11.19
- Keycloak IDP an Signing Server anbinden
- Unit Test fürs OIDC Karussell
- Kleine Fixes

#### 29.11.19
- Arbeiten am Verifier
    - add offline jwt verifizerung
    - eigene id tokens für jwt verifizerung
    - beginn signatur verifizierung
    - alle verifier zusammenhängen
- Arbeiten am Signing Server
    - CSR Fixen (subjectAltName ist immer eine Liste)
    - CA anbinden
    - CSR zur CA senden und Zertifikat erhalten
    - Protobuf anbinden
    - Signaturdatei beginnen zu erstellen
    - TSA Anbindung beginnen
    
#### 30.11.19
- Arbeiten am Signing Server
    - TSA Anbindung fertigstellen
    - CRL & OCSP holen
    - PKCS7 bauen
    - Über Bouncycastle fluchen
- Arbeiten am Verifier
    - Verbesserungen für concurrent verifying
- CA fixes für korrektes CA bundle
    
#### 01.12.19
- Arbeiten am Signing Server
    - CRL & OCSP in PKCS7 einbauen
    - CA Cert Bundle in PKCS7 einbauen
    - Erster Versuch eine Signaturdatei zu bauen
    - TSA Anfrage anpassen so dass Certs in der Antwort mitgeliefert werden
    - Signaturdateiformat wesentlich vereinfachen
    - Multisignaturen wesentlich vereinfachen
    - Soweit möglich http Anfragen parallelisieren (coroutines)
- Arbeiten am Verifier
    - anpassungen auf neues signaturformat
    
#### 02.12.19
- Verbessern Dokumentation des Login-Prozesses
- Fehler beheben in der Signing Server Implementation

#### 03.12.19
- Beginn Arbeit am Frontend
    - UI Design (Bootstrap CSS)
    - File Input entgegennehmen
    - Integration WASM Hashing Komponenten
    - Hashing Queue für sequentielles single-threaded hashing
    
#### 04.12.19
- Frontend:
    - IDP Karussell
    - OIDC Callback Zielseite
    - Daten in localstorage persistieren
    
