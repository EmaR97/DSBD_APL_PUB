# **Sistema di Monitoraggio delle Telecamere Distribuite**

### **Indice**

1. [Introduzione](#introduzione)
2. [Controllore di Base della Telecamera](#controllore-di-base-della-telecamera)
    - 2.1 [Caratteristiche Principali](#caratteristiche-principali)
    - 2.2 [Potenzialità per Sviluppi Futuri](#potenzialità-per-sviluppi-futuri)
3. [Server di Elaborazione](#server-di-elaborazione)
    - 3.1 [Responsabilità Chiave](#responsabilità-chiave)
    - 3.2 [Deployment Dinamico dei Server](#deployment-dinamico-dei-server)
4. [Server di Autenticazione](#server-di-autenticazione)
    - 4.1 [Caratteristiche Principali](#caratteristiche-principali-1)
    - 4.2 [Misure di Sicurezza](#misure-di-sicurezza)
    - 4.3 [Integrazione con Altri Componenti](#integrazione-con-altri-componenti)
5. [Server Principale](#server-principale)
    * 5.1 [Responsabilità Chiave](#responsabilità-chiave-1)
    * 5.2 [Miglioramento dell'Esperienza Utente ed Efficienza del Sistema](#miglioramento-dellesperienza-utente-ed-efficienza-del-sistema)
6. [Server dei Comandi](#server-dei-comandi)
    - 6.1 [Responsabilità Chiave](#responsabilità-chiave-1)
    - 6.2 [Miglioramento della Affidabilità per le Telecamere Remote](#miglioramento-della-affidabilità-per-le-telecamere-remote)
    - 6.3 [Scalabilità e Flessibilità](#scalabilità-e-flessibilità)
7. [Servizio di Sottoscrizione alle Notifiche](#servizio-di-sottoscrizione-alle-notifiche)
    - 7.1 [Funzionalità Chiave](#funzionalità-chiave)
    - 7.2 [Miglioramento del Controllo e della Personalizzazione Utente](#miglioramento-del-controllo-e-della-personalizzazione-utente)
8. [Servizio di Notifiche](#servizio-di-notifiche)
    - 8.1 [Funzionalità Chiave](#funzionalità-chiave-1)
9. [Sla Manager](#sla-manager)
    - 9.1 [Funzionalità Chiave](#funzionalità-chiave-2)
    - 9.2 [Aggiornamento Continuo del Modello](#aggiornamento-continuo-del-modello)
10. [Archiviazione Coerente dei Dati con MongoDB](#archiviazione-coerente-dei-dati-con-mongodb)
11. [Load Balancing e Routing in Docker e K8s](#load-balancing-e-routing-in-docker-e-k8s)
12. [Url pre-firmati per scaricare da Minio](#url-pre-firmati-per-scaricare-da-minio)
13. [Conclusioni e Sviluppi Futuri](#conclusioni-e-sviluppi-futuri)
14. [Istruzioni per Build e Deploy](Resource/build%20&%20deploy.md)
   

### **Introduzione:**

Il Sistema di Monitoraggio delle Telecamere Distribuite rappresenta una soluzione per soddisfare le crescenti esigenze di sorveglianza e rilevamento degli eventi in contesti distribuiti. Fondato su una serie di componenti come Kafka, MongoDB e Kubernetes, questo sistema si distingue per la sua capacità di catturare, elaborare e notificare eventi in tempo reale, mantenendo un'elevata scalabilità, efficienza e sicurezza.

Al centro di questo ecosistema si trovano le telecamere dislocate in diverse posizioni, responsabili della cattura dei frame che vengono successivamente inviati ai server di elaborazione attraverso il sistema di messaggistica Kafka. Questi server, utilizzando algoritmi di riconoscimento dei pedoni, elaborano i frame e li archiviano in MinIO, creando così un repository centralizzato per le immagini gestite. Il Server Principale, a sua volta, gestisce le informazioni e i frame delle telecamere, mentre il Server di Autenticazione garantisce un accesso sicuro agli utenti e alle loro telecamere associate.

Attraverso l'implementazione di un bot Telegram, il servizio di notifica consente agli utenti di interagire direttamente con il sistema tramite la piattaforma Telegram. Gli utenti possono gestire le proprie preferenze di notifica, decidendo quali eventi desiderano monitorare e quali tipologie di avvisi desiderano ricevere, tramite telegram stesso.

L'utilizzo GRPC, formato Proto e API-Gateway consente una comunicazione sicura e strutturata tra i vari componenti del sistema. Inoltre, l'architettura del sistema enfatizza la coerenza dei dati attraverso una gestione escusiva dei dati specifici ai vari servizi e l'accesso sicuro tramite API-Gateway, fornendo agli utenti un accesso diretto ai frame elaborati tramite URL pre-firmati.

Guardando al futuro, il sistema potrebbe beneficiare di ulteriori ottimizzazioni per migliorare la reattività in tempo reale, implementando meccanismi per il scaling dinamico e perfezionando l'interfaccia utente per offrire un'esperienza più ricca e fluida agli utenti. Inoltre, esplorare opportunità di integrazione con sistemi esterni potrebbe potenziare ulteriormente le capacità di riconoscimento e l'intelligenza complessiva del sistema, rendendolo ancora più versatile e adattabile alle esigenze emergenti nel campo della sorveglianza e del rilevamento degli eventi.

### Schema di Componenti e Comunicazione 
![Schema di Comunicazione](Resource/communication_scheme.jpg)


---

### **Controllore di Base della Telecamera:**

**[Src](CPP/src/cam_controller.cpp)**

Il Controllore di Base della Telecamera si occupa di rappresentare e implementare le funzionalità di base fornite dalla telecamera nel sistema e supporta eventuali estensioni per funzionalità avanzate, grazie alla sua struttura modulare.

#### *Caratteristiche Principali:*

1. **Distribuzione di Frame tramite Kafka:**

    - Il controllore cattura i frame e li distribuisce in modo efficiente ai servers di elaborazione attraverso un topic Kafka dedicato.
    - Utilizza Kafka per garantire una consegna veloce dei frame con tolleranza agli errori.

2. **Elaborazione di Comandi tramite RabbitMQ e MQTT:**

    - Ascolta i comandi tramite RabbitMQ utilizzando il protocollo MQTT.
    - Esegue compiti basati sui comandi ricevuti, consentendo un controllo dinamico del sistema della telecamera.

3. **Standardizzazione dei Messaggi:**

    - I messaggi in ingresso e in uscita seguono il formato Proto standard di Google.
    - Questo assicura un protocollo di comunicazione coerente, facilitando l'interoperabilità e l'integrazione con altri componenti.

#### *Potenzialità per Sviluppi Futuri:*

- Il design modulare del Controllore di Base della Telecamera consente un'integrazione senza soluzione di continuità di
  funzionalità avanzate. Gli sviluppatori possono utilizzare questo codice come base per implementare ulteriori
  funzioni, migliorando le capacità del sistema di monitoraggio delle telecamere.

<!--In sintesi, il Controllore di Base della Telecamera offre un solido punto di partenza per il sistema di monitoraggio
delle telecamere distribuite, fornendo funzionalità essenziali di distribuzione di frame ed elaborazione di comandi in
un formato standardizzato. -->

---

### **Server di Elaborazione:**

**[Src](CPP/src/processing_server.cpp)**

I server di elaborazione si concentrano sull'applicazione efficace dell'algoritmo di riconoscimento dei pedoni e sulla gestione successiva delle immagini elaborate.

#### *Responsabilità Chiave:*

1. **Elaborazione Parallelizzata:**

    - Diversi server di elaborazione si connettono al medesimo gruppo Kafka su un topic dedicato, permettendo la distribuzione del carico di lavoro per l'applicazione dell'algoritmo di riconoscimento dei pedoni.
    - Questo approccio parallelo migliora la scalabilità e la reattività del sistema, contribuendo alla sua efficienza complessiva.

2. **Elaborazione Algoritmica e Archiviazione delle Immagini:**

    - I server di elaborazione applicano l'algoritmo di riconoscimento dei pedoni ai frame ricevuti per un'analisi accurata e tempestiva.
    - Le immagini elaborate vengono trasferite a MinIO per l'archiviazione, creando un repository centralizzato per i dati storici.

3. **Segnalazione dei Risultati al Server Principale:**

    - I risultati dell'algoritmo di riconoscimento dei pedoni vengono trasmessi al server principale per l'archiviazione e ulteriori elaborazioni.
    - Ciò agevola un'analisi completa e supporta le decisioni basate sulle informazioni rilevate.

4. **Metriche di Utilizzo delle Risorse con Prometheus:**

    - Prometheus è integrato per raccogliere metriche essenziali al fine di ottimizzare l'utilizzo delle risorse:

        - Il parametro **time_since_message_creation
          ** valuta il tempo trascorso dalla creazione del frame al suo arrivo al server di elaborazione, fornendo informazioni sulla latenza causata da una capacità di elaborazione insufficiente.
        - Il parametro **working_time
          ** valuta il rapporto tra il tempo dedicato al lavoro sui frame e il tempo inattivo dei server di elaborazione, contribuendo alle decisioni di allocazione delle risorse e consentendo di individuare riduzioni nel carico di lavoro.

#### *Deployment Dinamico dei Server:*

- Le metriche raccolte, insieme a quelle ottenute tramite un exporter posizionato sul nodo di Kafka, fungono da base per valutare la capacità di elaborazione delle immagini attuale del sistema. Queste metriche potenzialmente supportano lo sviluppo di un deployment dinamico dei server, regolando il numero di server di elaborazione deployati in modo dinamico per garantire un utilizzo ottimale delle risorse e la reattività alle variazioni del carico di lavoro.

<!-- In sintesi, i server di elaborazione contribuiscono alla natura distribuita del sistema attraverso l'elaborazione
parallelizzata, l'ottimizzazione delle risorse mediante metriche di Prometheus e la scalabilità dinamica.-->

---

### **Server di Autenticazione:**

**[Src](GO/src/server_auth/main.go)**

Il Server di Autenticazione fornisce un meccanismo di accesso sicuro e controllato per gli utenti e le loro telecamere associate, alle funzionalità offerte dal sistema.

#### *Caratteristiche Principali:*

1. **Integrazione con MongoDB per la Gestione delle Credenziali Utente:**

    - Il Server di Autenticazione si connette a un server MongoDB per la gestione delle credenziali utente, compresi
      dettagli di autenticazione per utenti e relative telecamere.
    - Concedere privilegi di accesso in base ai ruoli e alle autorizzazioni degli utenti assicura un accesso controllato
      alle funzionalità del sistema.

2. **Autenticazione Utente e Assegnazione delle Credenziali:**

    - Al login avvenuto con successo, il Server di Autenticazione genera una password per la telecamera per l'utente,
      utilizzata dalle telecamere associate per accedere al sistema.
    - Un token temporaneo viene fornito al login, fungendo da identificatore per richieste successive al sistema.

3. **Verifica del Token per l'Autorizzazione delle Richieste:**

    - Ogni richiesta di sistema include un token, verificato dal Server di Autenticazione. L'accesso alla funzionalità
      richiesta è concesso solo al successo della conferma del token, garantendo un ambiente sicuro.
    - L'uso di token temporanei aumenta la sicurezza mediante l'aggiornamento regolare dell'identificazione dell'utente,
      riducendo il rischio di accessi non autorizzati.

4. **Provision delle Credenziali alle Telecamere:**

    - Le telecamere, al login avvenuto con successo, ricevono le credenziali necessarie per accedere ai servizi RabbitMQ
      e Kafka, garantendo un'integrazione senza soluzione di continuità nel sistema distribuito per la comunicazione e
      lo scambio di dati in tempo reale.


#### API Implementate
* **GET _/access/login_**: Restituisce la pagina di login per l'autenticazione dell'utente.
* **POST _/access/login_**: Gestisce l'autenticazione dell'utente. Richiede i parametri `username`, `password`, e `cam_id`.
* **GET _/access/signup_**: Restituisce la pagina di registrazione per la creazione di un nuovo account utente.
* **POST _/access/signup_**: Gestisce la registrazione di un nuovo account utente. Richiede i parametri `username`, `email`, e `password`.
* **GET _/access/logout_**: Gestisce il logout dell'utente.
* **POST _/access/verify_**: Verifica la validità di un token. Richiede il parametro `token`.
#### *Misure di Sicurezza:*

- Il Server di Autenticazione agisce come un guardiano, verificando la legittimità delle richieste e garantendo che solo
  utenti e telecamere autorizzati possano interagire con il sistema.

#### *Integrazione con Altri Componenti:*

- Il Server di Autenticazione svolge un ruolo vitale nell'orchestrare la comunicazione sicura e l'interazione tra
  utenti, telecamere e servizi di sistema, inclusi RabbitMQ e Kafka.

<!-- In sintesi, il Server di Autenticazione stabilisce un robusto framework di autenticazione e autorizzazione, integrandosi
con MongoDB per la gestione delle credenziali e garantendo una comunicazione sicura all'interno del sistema di
monitoraggio delle telecamere distribuite. -->

---

### **Server Principale:**

**[Src](GO/src/server_main/main.go)**

Il Server Principale supervisiona
la gestione delle telecamere, le registrazioni degli utenti e l'archiviazione di frame e informazioni pertinenti.

#### *Responsabilità Chiave:*

1. **Gestione e Registrazione delle Telecamere:**

    - Gli utenti possono registrare nuove telecamere effettuando richieste al Server Principale. L'ID della telecamera
      generato viene successivamente utilizzato per il login della telecamera, agevolando l'integrazione senza soluzione
      di continuità dei nuovi dispositivi.
    - Questo processo di registrazione consente agli utenti di espandere senza sforzo la propria rete di telecamere.

2. **Archiviazione di Informazioni su Frame e Telecamera:**

    - Il Server Principale è responsabile della gestione delle informazioni delle telecamere e dell'archiviazione sicura
      dei frame in MongoDB.
    - I dettagli della telecamera, inclusi l'identificazione e le credenziali di accesso, vengono memorizzati per un
      accesso e una gestione efficienti.

3. **Integrazione con Kafka per Informazioni sui Frame Elaborati:**

    - Il Server Principale è registrato su un topic Kafka dove i server di elaborazione depositano informazioni relative
      ai frame elaborati.
    - Questa integrazione consente al Server Principale di raccogliere dati cruciali sull'elaborazione dei frame,
      migliorando la comprensione complessiva dell'identificazione di MinIO, dettagli aggiuntivi dell'immagine e l'esito
      del riconoscimento dei pedoni.

4. **Accesso degli Utenti ai Flussi Video delle Telecamere:**

    - Attraverso il Server Principale, gli utenti possono accedere ai flussi video delle loro telecamere registrate, fornendo capacità di monitoraggio in tempo reale.
    - Un url pre-firmato viene generato per permettere l'accesso diretto dell'utente a Minio

5. **Recupero di Informazioni per Altri Servizi:**

    - Altri servizi all'interno del sistema possono interrogare il Server Principale per accedere a informazioni sulle
      telecamere registrate, garantendo coerenza e affidabilità nell'accesso ai dettagli delle telecamere.

6. **Notifiche Positive per il Riconoscimento dei Pedoni:**

    - Alla ricezione di messaggi relativi alle immagini elaborate, se un'immagine viene positivamente riconosciuta per i
      pedoni, il Server Principale invia un messaggio sul topic di notifica su Kafka.
    - Questo messaggio viene quindi elaborato dal Servizio di Notifiche, consentendo agli utenti di ricevere avvisi
      tempestivi e notifiche sull'attività dei pedoni rilevata.

7. **Gestione dell'Eliminazione Automatica di Dati Obsoleti:**

    - Il sistema implementa un meccanismo di pulizia automatica per le immagini salvate.

    - Dopo un tempo predefinito, le immagini elaborate archiviate su MinIO e le corrispondenti informazioni su MongoDB
      vengono eliminate per garantire l'ottimizzazione dello spazio di archiviazione e la gestione efficiente delle
      risorse del sistema.

    - Questa pratica assicura che solo dati pertinenti e recenti siano conservati nel sistema, riducendo l'ingombro e
      contribuendo alla performance ottimale del sistema nel lungo termine.

#### API Implementate
* **GET _/api/camera/:id/:lastSeen_**: Restituisce i frame successivi al timestamp specificato per una determinata fotocamera. `:id` rappresenta l'identificatore univoco della fotocamera, mentre `:lastSeen` indica il timestamp dell'ultimo frame visualizzato.
* **POST _/api/camera_**: Crea una nuova fotocamera. Non richiede parametri aggiuntivi.
* **POST _/api/camera/login_**: Effettua il login per una fotocamera. Non richiede parametri aggiuntivi.
* **GET _/api/videoFeed/:id_**: Permette di visualizzare il feed video per una specifica fotocamera. `:id` è l'identificatore univoco della fotocamera di cui si desidera visualizzare il feed.
* **GET _/api/videoFeed_**: Restituisce l'elenco delle telecamere disponibili. Non richiede parametri aggiuntivi.
#### *Miglioramento dell'Esperienza Utente ed Efficienza del Sistema:*

- Il Server Principale funge da perno, fornendo un'interfaccia coesa per gli utenti per gestire le telecamere, accedere
  ai flussi video e ricevere notifiche.
- La sua integrazione con Kafka migliora l'efficienza e la reattività del sistema.

<!-- In sintesi, il Server Principale svolge un ruolo cruciale nella gestione delle telecamere, nelle interazioni degli
utenti e nel flusso senza soluzione di continuità delle informazioni all'interno del sistema distribuito di monitoraggio
delle telecamere. -->

---

### **Server dei Comandi:**

**[Src](GO/src/server_command/main.go)**

Il Server dei Comandi consente agli utenti di inviare comandi tramite richieste API, i quali vengono trasferiti in modo efficiente alle telecamere specificate, affrontando le sfide delle connessioni remote delle telecamere.

#### *Responsabilità Chiave:*

1. **Ricezione delle Richieste API:**
    - Il Server dei Comandi agisce come intermediario tra l'interfaccia utente e le telecamere nel sistema, accogliendo le richieste API dagli utenti.

2. **Formattazione e Standardizzazione dei Comandi:**
    - I comandi ricevuti vengono formattati e standardizzati utilizzando il formato Proto, garantendo così un protocollo di comunicazione uniforme e strutturato.
    - Questo processo di standardizzazione migliora l'interoperabilità e semplifica l'integrazione con diversi componenti del sistema.

3. **Consegna Affidabile dei Messaggi tramite MQTT:**
    - MQTT viene impiegato come protocollo di comunicazione per assicurare una consegna affidabile dei messaggi, particolarmente vantaggiosa in scenari con connessioni instabili, come quelle tipiche delle telecamere remote degli utenti.
    - La conferma di consegna del messaggio garantita da MQTT assicura che i comandi raggiungano le telecamere specificate anche in condizioni di rete difficili.

4. **Comunicazione con RabbitMQ:**
    - I comandi formattati vengono trasmessi al topic della telecamera specificata all'interno di RabbitMQ tramite MQTT, stabilendo un canale di comunicazione affidabile ed efficiente tra il Server dei Comandi e le telecamere.

#### API Implementate
* **POST _/commands/:id_**: Invia un comando alla fotocamera identificata da `:id`.

#### *Miglioramento della Affidabilità per le Telecamere Remote:*

- La scelta di MQTT come protocollo di comunicazione offre una consegna affidabile dell'ultimo messaggio, sovrascritto dall'arrivo del successivo.

#### *Scalabilità e Flessibilità:*

- Il design del Server dei Comandi è progettato per consentire la scalabilità, gestendo un numero crescente di utenti e telecamere. Inoltre, grazie all'utilizzo di un IDL come Proto, garantisce un'integrazione fluida con varie interfacce fornite dalle telecamere.

<!--In sintesi, il Server dei Comandi svolge un ruolo cruciale nel facilitare i comandi degli utenti, standardizzandoli e garantendo una consegna affidabile alle telecamere remote tramite l'infrastruttura robusta di MQTT e RabbitMQ.-->
---

### **Servizio di Sottoscrizione alle Notifiche:**

**[Src](Python/src/conversation_bot/main.py)**

Il Servizio di Sottoscrizione Notifiche è un bot Telegram che intrattiene conversazioni con gli utenti sulla piattaforma
Telegram, offrendo un'interfaccia per gestire le preferenze e le sottoscrizioni alle
notifiche degli utenti.

#### *Funzionalità Chiave:*

1. **Bot Telegram per l'Interazione con l'Utente:**

    - Operando come un bot Telegram, il servizio consente agli utenti di accedere direttamente alle sue funzionalità
      attraverso la piattaforma Telegram.
    - Gli utenti si autenticano con le proprie credenziali per stabilire una connessione sicura.

2. **Archiviazione dell'ID Utente:**

    - Il servizio memorizza l'ID utente di Telegram per mantenere un registro della conversazione e delle preferenze
      dell'utente, fungendo da identificatore chiave per associare gli utenti alle loro preferenze di notifica.

3. **Gestione delle Sottoscrizioni:**

    - Gli utenti possono gestire le preferenze di notifica, incluso l'abbonamento alle telecamere di loro proprietà.
    - Dettagli di sottoscrizione, come l'intervallo di tempo tra le notifiche e l'orario preferito per ricevere
      notifiche, possono essere specificati dall'utente.

4. **Funzione di Annullamento dell'Abbonamento:**

    - Gli utenti possono annullare l'abbonamento per non ricevere notifiche per una telecamera specifica, offrendo
      flessibilità e assicurando che gli utenti ricevano solo notifiche rilevanti.

5. **Interfaccia GRPC per il Recupero delle Informazioni:**

    - Il Servizio di Sottoscrizione alle Notifiche espone un'interfaccia GRPC per consentire ad altri componenti, come il
      Servizio di Notifiche, di ottenere le informazioni necessarie sulle sottoscrizioni degli utenti.

#### *Miglioramento del Controllo e della Personalizzazione Utente:*

- Il servizio dà potere agli utenti permettendo loro di personalizzare le preferenze di notifica, specificando le
  telecamere di interesse, l'orario delle notifiche e la possibilità di optare per l'uscita in qualsiasi momento.



<!-- In sintesi, il Servizio di Sottoscrizione Notifiche funge da interfaccia utente amichevole sulla piattaforma Telegram,
consentendo agli utenti di gestire le loro preferenze di notifica e interagire con il più ampio sistema di notifiche.-->
---

### **Servizio di Notifiche:**

**[Src](Python/src/notification_bot/main.py)**

Il Servizio di Notifiche è un bot Telegram responsabile del consumo di messaggi dal topic di notifica Kafka e della
notifica efficiente degli utenti che hanno manifestato interesse in specifiche notifiche.

#### *Funzionalità Chiave:*

1. **Consumo di Messaggi da Kafka:**

    - Il Servizio di Notifiche consuma continuamente i messaggi dal topic di notifica Kafka.
    - Questi messaggi contengono informazioni su un riconoscimento positivo di pedoni.

2. **Recupero delle Informazioni sugli Utenti:**

    - Con ogni messaggio consumato, il Servizio di Notifiche richiede informazioni al Servizio di Sottoscrizione alle
      Notifiche tramite l'interfaccia GRPC.
    - Questa richiesta aiuta a identificare tutti gli utenti interessati a ricevere notifiche relative all'evento
      specifico menzionato nel messaggio Kafka.

3. **Notifiche Telegram:**

    - Il servizio utilizza gli ID di Telegram registrati ottenuti dal Servizio di Sottoscrizione Notifiche per
      notificare gli utenti sull'evento.
    - Le notifiche vengono inviate direttamente agli utenti sulla piattaforma Telegram, fornendo avvisi in tempo reale
      su attività di riconoscimento dei pedoni.

4. **Mirato Efficientemente agli Utenti:**

    - Sfruttando le informazioni ottenute dal Servizio di Sottoscrizione Notifiche, il Servizio di Notifiche garantisce
      che le notifiche siano indirizzate solo agli utenti che hanno manifestato interesse in eventi specifici della
      telecamera.

5. **Scalabilità e Reattività:**

    - Il design del Servizio di Notifiche supporta la scalabilità, gestendo efficientemente un numero crescente di
      notifiche e utenti.
    - Il servizio risponde prontamente ai messaggi Kafka in arrivo, garantendo notifiche tempestive agli utenti
      interessati.

<!--In sintesi, il Servizio di Notifiche svolge un ruolo cruciale nell'ultimo passo del processo di notifica, garantendo che
gli utenti che si sono abbonati a eventi specifici ricevano avvisi tempestivi e personalizzati sulla piattaforma
Telegram.-->

---

### **Sla Manager**

**[Src](Python/src/sla_manager/main.py)**

Il server SlaManager è un componente cruciale del sistema di monitoraggio delle telecamere distribuite. Consente l'aggiornamento dinamico degli SLA e utilizza tecniche di analisi dei dati per valutare la probabilità di violazioni delle SLA definite e per adattare i modelli alle condizioni correnti dei dati.

#### *Funzionalità Chiave:*

1. **Stima della Probabilità di Violazioni:**

    - Utilizza la distribuzione gaussiana per stimare la probabilità che una variabile metrica superi determinati limiti, basandosi su tendenze storiche e parametri del modello.
    - Calcola la probabilità di violazione degli SLA all'interno di specifici intervalli di tempo, consentendo la previsione e la gestione pro-attiva dei problemi.

2. **Adattamento del Modello:**

    - Implementa tecniche di adattamento del modello per incorporare nuovi dati e modifiche nelle condizioni di sistema.
    - Utilizza algoritmi di fitting polinomiale, decomposizione stagionale e stima dell'errore per adattare i modelli alle variazioni nei dati di telecamere e garantire predizioni accurate.

3. **Aggiornamento degli SLA:**

    - Consente l'aggiornamento dinamico degli SLA (Service Level Agreement) in base alle esigenze operative e ai cambiamenti nelle prestazioni del sistema.
    - Incorpora nuovi requisiti degli SLA nel processo di stima della probabilità e nell'adattamento del modello per garantire la conformità allo standard di servizio.

#### *Aggiornamento Continuo del Modello:*

- Il server SlaManager si aggiorna costantemente con nuovi dati e informazioni sugli SLA, garantendo una precisione e una affidabilità continue nella stima delle probabilità e nel fitting dei modelli.

<!--In sintesi, il Server di Stima della Probabilità Gaussiana e Adattamento del Modello fornisce un'analisi sofisticata delle metriche di telecamere, consentendo la valutazione delle violazioni degli SLA e l'adattamento dinamico dei modelli per rispondere alle condizioni del sistema.-->

---

### **Archiviazione Coerente dei Dati con MongoDB:**

MongoDB funge da database centralizzato per archiviare varie categorie di dati, mantenendo un approccio
strutturato e organizzato alla gestione dei dati.

Per mantenere la coerenza dei dati, ogni categoria di dati è accessibile attraverso un componente singolare,
minimizzando il rischio di inconsistenze dei dati e garantendo interazioni ben definite con tipi di dati
specifici.

---

### **Load Balancing e Routing in Docker e K8s:**

Nel nostro ambiente Docker e Kubernetes (K8s), il load balancing e il routing sono gestiti attraverso l'utilizzo di Nginx o Ingress. Questa configurazione consente un accesso strutturato tramite API alle funzionalità del sistema, semplificando notevolmente l'implementazione del load balancing e del routing. Ciò assicura una distribuzione uniforme del traffico e un re-indirizzamento efficiente delle richieste.

In futuro, potremmo considerare l'implementazione di funzionalità più avanzate utilizzando i Gateway API in Kubernetes (K8s). Questa tecnologia offre una maggiore personalizzazione e una gamma più ampia di funzionalità rispetto ai gateway tradizionali come Nginx. In particolare, consente la gestione non solo della comunicazione tramite HTTP, ma anche di protocolli come gRPC e TCP.

---

### **Url pre-firmati per scaricare da Minio:**

Il Server Principale facilita l'accesso diretto ai frame elaborati archiviati in MinIO generando URL pre-firmati. Questo approccio minimizza la presenza di intermediari non necessari, ottimizzando la velocità e fornendo agli utenti un accesso efficiente ai dati archiviati. In questo modo, viene garantito un accesso diretto, ma controllato, al sistema di archiviazione di Minio.

---

### **Conclusioni e Sviluppi Futuri:**

Il sistema di monitoraggio delle telecamere presenta alcune aree che potrebbero essere ulteriormente ottimizzate per migliorare le prestazioni e l'esperienza utente. In
particolare:

#### **Reattività in Tempo Reale:**

- Nonostante l'efficienza complessiva del sistema, potrebbero essere esplorate ulteriori ottimizzazioni per migliorare
  la reattività in tempo reale, specialmente in scenari con carichi di lavoro variabili.

#### **Scaling Dinamico:**

- Il sistema potrebbe beneficiare di meccanismi per il scaling dinamico, regolando automaticamente il
  numero di server di elaborazione in base a metriche in tempo reale per garantire un utilizzo ottimale delle risorse.

#### **Miglioramenti dell'Interfaccia Utente:**

- L'interfaccia utente, specialmente nei bot Telegram, potrebbe essere perfezionata per offrire più funzionalità e una
  maggiore fluidità nell'esperienza utente, integrando contenuti multimediali o comandi aggiuntivi per l'interazione con le telecamere.

#### **Integrazione con Sistemi Esterni:**

- Esplorare possibilità di integrazione del sistema con servizi esterni o framework di intelligenza artificiale potrebbe
  potenziare ulteriormente le capacità di riconoscimento dei pedoni e l'intelligenza complessiva del sistema.
---

