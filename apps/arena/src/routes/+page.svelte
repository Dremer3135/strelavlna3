<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import Main from './Main.svelte';
  import { probs } from '$lib/stores/probs';
  import { currentState } from '$lib/stores/state';
  import type { MessageType } from '$lib/types';
  
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
          name: "tohle je jmeno",
          text: "tohle jetextulohy",
          answer: "",
          diff: "A",
          images: [],
          focusedBy: ["lalala", "bablbamId"],
          chat: [{
            origin: "recieved",
            type: "message",
            value: "Hello my friend!"
          },
          {
            origin: "sent",
            type: "message",
            value: "Hi!"
          },
        ]
        },
        "ijfs": {
          id: "ijfs",
          name: "tohle je jmeno",
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
        
      }
    });
    currentState.update(state => {
      return {
        teamName: "Bambuláci 4. trida",
        money: 1534,
        myId: "bablbamId"
      }
    })
  });

  function handleBuy(diff: string) {
    console.log("Buying:", diff);
  }
  function handleChat(probId: string, message: MessageType) {
    console.log("Message of type '" + message.type + "': " + message.value);
  }
  function handleSell(probId: string) {
    console.log("Selling:", probId);
  }
  function handleFocus(probId: string) {
    console.log("Focusing:", probId);
  }
</script>

<main>
  <Main buy={handleBuy} chat={handleChat} sell={handleSell} focus={handleFocus} />
</main>

<style lang="scss">
  main {
    display: flex;
    flex-direction: column;
    width: 100vw;
    height: 100vh;
  }
</style>