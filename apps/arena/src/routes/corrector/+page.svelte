<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import Main from './Main.svelte';
  import { probs } from '$lib/stores/probs';
  import { currentState, wsConnected } from '$lib/stores/state';
  import type { MessageType } from '$lib/types';
    import Disconnected from '$lib/assets/components/Disconnected.svelte';
  
  onMount(() => {

    // const socket = new WebSocket(`https://sv.skrat.org/ws?token=${$page.data.token}`)
    // socket.addEventListener("message", (event) => {
    //   console.log(typeof event.data)
    // });


    probs.update(_ => {
      return { 
        "askhdj": {
          id: "askhdj",
          name: "Tohle je jmeno",
          text: "tohle jetextulohy",
          answer: "",
          diff: "A",
          images: [],
          focusedBy: ["s", "lalala"],
          chat: []
        },
        "asjldh": {
          id: "asjldh",
          name: "tohle je jmeno lalal",
          text: "tohle jetextulohy",
          answer: "",
          diff: "A",
          images: [],
          focusedBy: ["lalala", "bablbamId"],
          chat: [
            {
              origin: "recieved",
              type: "message",
              value: "Hello my friend!",
              sentTime: new Date('2025-11-27T14:05:30Z')
            },
            {
              origin: "sent",
              type: "message",
              value: "Hi!",
              sentTime: new Date('2025-11-27T14:05:30Z')
            },
            {
              origin: "sent",
              type: "answer",
              value: "Odpoved je lalala",
              sentTime: new Date('2025-11-27T14:05:30Z')
            },
            {
              origin: "sent",
              type: "grade",
              value: "incorrect",
              sentTime: new Date('2025-11-27T14:05:30Z')
            },
            {
              origin: "recieved",
              type: "message",
              value: "prosiiim",
              sentTime: new Date('2025-11-27T14:05:30Z')
            },
            {
              origin: "recieved",
              type: "answer",
              value: "Odpoved neni lalala",
              sentTime: new Date('2025-11-27T14:05:30Z')
            },
            {
              origin: "recieved",
              type: "copy",
              value: "tohlrje je neco co dam chatgpt",
              sentTime: new Date('2025-11-27T14:05:30Z')
            },
            {
              origin: "recieved",
              type: "window-focus",
              value: "blur honzik",
              sentTime: new Date('2025-11-27T14:05:30Z')
            },
            {
              origin: "recieved",
              type: "window-focus",
              value: "focus honzik",
              sentTime: new Date('2025-11-27T14:05:30Z')
            },
            {
              origin: "recieved",
              type: "paste",
              value: "tohlrje je z chatgpt",
              sentTime: new Date('2025-11-27T14:05:30Z')
            },
            {
              origin: "recieved",
              type: "paste",
              value: "tohlrje je z chatgpt",
              sentTime: new Date('2025-11-27T14:05:30Z')
            },
            {
              origin: "recieved",
              type: "paste",
              value: "tohlrje je z chatgpt",
              sentTime: new Date('2025-11-27T14:05:30Z')
            },
            {
              origin: "recieved",
              type: "paste",
              value: "tohlrje je z chatgpt",
              sentTime: new Date('2025-11-27T14:05:30Z')
            },

            // {
            //   origin: "recieved",
            //   type: "grade",
            //   value: "correct",
            //   sentTime: new Date('2025-11-27T14:05:30Z')
            // },
          ]
        },
        "ijfs": {
          id: "ijfs",
          name: "tohle je jmenoslav",
          text: "tohle jetextulohy",
          answer: "",
          diff: "A",
          images: [],
          focusedBy: [],
          chat: []
        },
        "alsAL": {
          id: "alsAL",
          name: "tohle je jmeno",
          text: "tohle jetextulohy",
          answer: "",
          diff: "A",
          images: [],
          focusedBy: ["adf", "ads"],
          chat: []
        },
        "PIJFSLKHBKJ": {
          id: "PIJFSLKHBKJ",
          name: "tohle je jmeno",
          text: "tohle jetextulohy",
          answer: "",
          diff: "A",
          images: [],
          focusedBy: [],
          chat: []
        },
        "baa": {
          id: "PIJFSLKHBKJ",
          name: "tohle je jmeno",
          text: "tohle jetextulohy",
          answer: "",
          diff: "A",
          images: [],
          focusedBy: [],
          chat: []
        },
        "a": {
          id: "PIJFSLKHBKJ",
          name: "tohle je jmeno",
          text: "tohle jetextulohy",
          answer: "",
          diff: "A",
          images: [],
          focusedBy: [],
          chat: []
        },
        "b": {
          id: "PIJFSLKHBKJ",
          name: "tohle je jmeno",
          text: "tohle jetextulohy",
          answer: "",
          diff: "A",
          images: [],
          focusedBy: [],
          chat: []
        },
        "c": {
          id: "PIJFSLKHBKJ",
          name: "tohle je jmeno",
          text: "tohle jetextulohy",
          answer: "",
          diff: "A",
          images: [],
          focusedBy: [],
          chat: []
        },
        "d": {
          id: "PIJFSLKHBKJ",
          name: "tohle je jmeno",
          text: "tohle jetextulohy",
          answer: "",
          diff: "A",
          images: [],
          focusedBy: [],
          chat: []
        },
        
      }
    });
    currentState.update(state => {
      return {
        teamName: "Bambuláci 4. trida",
        money: 80,
        myId: "bablbamId",
        probsRemaining: [10, 2, -1],
        pricesBuy: [10, 30, 80],
        pricesSell: [10, 15, 40],
        procesSolve: [15, 50, 200],
        start: new Date(Date.now() + 0.2*60*1000),
        end: new Date(Date.now() + 20*60*1000)
      }
    })
  });

  function handleChat(probId: string, message: Omit<MessageType, 'sentTime'>) {
    console.log("Message of type '" + message.type + "': " + message.value);
  }
  function handleFocus(probId: string) {
    console.log("Focusing:", probId);
  }
</script>

<main>
  {#if $wsConnected}
    <Main chat={handleChat} focus={handleFocus} />
  {:else}
    <Disconnected/>
  {/if}
</main>

<style lang="scss">
  main {
    flex-grow: 1;
    min-height: 0px;
    display: flex;
    flex-direction: column;
    width: 100%;
    height: 100%;
  }
</style>