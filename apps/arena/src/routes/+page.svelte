<script lang="ts">
  import { onMount } from 'svelte';
  import Main from './Main.svelte';
  import { probs } from '$lib/stores/probs';
  import { currentState, wsConnected } from '$lib/stores/state';
  import type { MessageType, CurrentState } from '$lib/types';
  import Waitroom from './Waitroom.svelte';
    import Disconnected from '$lib/assets/components/Disconnected.svelte';

  let { data } = $props();

  let socket: WebSocket;
  
  onMount(() => {
    socket = new WebSocket(`https://sv.skrat.org/ws?token=${data.token}`);

    socket.onopen = () => {
      wsConnected.set(true);
    };
    // socket.onclose = () => {
    //   wsConnected.set(false);
    // }
    // socket.onerror = () => {
    //   wsConnected.set(false);
    // }

    socket.addEventListener("message", (event) => {
      let data = JSON.parse(event.data);
      console.log(data);
      if (data.name === "initload") {
        currentState.set({
          teamName: data.data.teamname,
          money: data.data.money,
          myId: data.data.playerid,
          probsRemaining: [data.data.remprobs.A, data.data.remprobs.B, data.data.remprobs.C],
          pricesBuy: [data.data.buycost.A, data.data.buycost.B, data.data.buycost.C],
          pricesSell: [data.data.sellcost.A, data.data.sellcost.B, data.data.sellcost.C],
          procesSolve: [data.data.solvecost.A, data.data.solvecost.B, data.data.solvecost.C],
          start: new Date(Date.now() + data.data.start),
          end: new Date(Date.now() + data.data.end),
          runningState: data.data.state,
          rank: -1,
        });

        let nprobs: any = {};

        for (let prob of data.data.bought) {
          nprobs[prob.id] = {
            id: prob.id,
            name: prob.name,
            text: prob.text,
            answer: prob.answer,
            diff: prob.diff,
            images: prob.images,
            focusedBy: [],
            chat: data.data.tlines[prob.id].map((e: any) => {return {
              origin: e.mside === "admin" ? "sent" : "recieved",
              type: e.mtype,
              value: e.msg,
              sentTime: e.time,
            }}),
            owned: "bought",
          };
        }

        for (let prob of data.data.solved) {
          nprobs[prob.id] = {
            id: prob.id,
            name: prob.name,
            text: prob.text,
            answer: prob.answer,
            diff: prob.diff,
            images: prob.images,
            focusedBy: [],
            chat: data.data.tlines[prob.id].map((e: any) => {return {
              origin: e.mside === "admin" ? "sent" : "recieved",
              type: e.mtype,
              value: e.msg,
              sentTime: e.time,
            }}),
            owned: "solved",
          };
        }

        for (let prob of data.data.sold) {
          nprobs[prob.id] = {
            id: prob.id,
            name: prob.name,
            text: prob.text,
            answer: prob.answer,
            diff: prob.diff,
            images: prob.images,
            focusedBy: [],
            chat: data.data.tlines[prob.id].map((e: any) => {return {
              origin: e.mside === "admin" ? "sent" : "recieved",
              type: e.mtype,
              value: e.msg,
              sentTime: e.time,
            }}),
            owned: "sold",
          };
        }

        probs.set(nprobs);
      } else if (data.name === "focus") {
        probs.update((e) => {
          Object.keys(e).forEach((k) => {
            if (k === data.probid) {
              e[k].focusedBy = e[k].focusedBy.filter((x: any) => x !== $currentState.myId);
            } else {
              e[k].focusedBy.push($currentState.myId);
            }
          })
          return e;
        });
      } else if (data.name === "written") {
        probs.update((e) => {
          Object.keys(e).forEach((k) => {
            if (k === data.probid) {
              e[k].chat.push({
                origin: data.type,
                type: data.solve ? "answer" : "message",
                value: data.text,
                sentTime: data.time,
              });
            }
          })
          return e;
        });
      } else if (data.name === "bought") {
        probs.update((e) => {
          e[data.prob.id] = {
            id: data.prob.id,
            name: data.prob.name,
            text: data.prob.text,
            answer: data.prob.answer,
            diff: data.prob.diff,
            images: data.prob.images,
            focusedBy: [],
            chat: [],
            owned: "bought",
          };
          return e;
        });
        currentState.update((e) => { return {...e, money: data.money, probsRemaining: [data.remprobs.A, data.remprobs.B, data.remprobs.C]}})
      } else if (data.name === "sold") {
        probs.update((e) => {
          e[data.probid].owned = "sold";
          currentState.update((e) => { return {...e, money: data.money, probsRemaining: [data.remprobs.A, data.remprobs.B, data.remprobs.C]}})
          return e;
        });
      } else if (data.name === "solved") {
        probs.update((e) => {
          e[data.probid].owned = "solved";
          currentState.update((e) => { return {...e, money: data.money, probsRemaining: [data.remprobs.A, data.remprobs.B, data.remprobs.C]}})
          return e;
        });
      } else if (data.name === "start") {
        currentState.update((e) => { return {...e, runningState: "running"}})
      }
    });


    // probs.update(_ => {
    //   return { 
    //     "askhdj": {
    //       id: "askhdj",
    //       name: "Tohle je jmeno",
    //       text: "tohle jetextulohy",
    //       answer: "",
    //       diff: "A",
    //       images: [],
    //       focusedBy: ["s", "lalala"],
    //       chat: []
    //     },
    //     "asjldh": {
    //       id: "asjldh",
    //       name: "tohle je jmeno",
    //       text: "tohle jetextulohy",
    //       answer: "",
    //       diff: "A",
    //       images: [],
    //       focusedBy: ["lalala", "bablbamId"],
    //       chat: [
    //         {
    //           origin: "recieved",
    //           type: "message",
    //           value: "Hello my friend!",
    //           sentTime: new Date('2025-11-27T14:05:30Z')
    //         },
    //         {
    //           origin: "sent",
    //           type: "message",
    //           value: "Hi!",
    //           sentTime: new Date('2025-11-27T14:05:30Z')
    //         },
    //         {
    //           origin: "sent",
    //           type: "answer",
    //           value: "Odpoved je lalala",
    //           sentTime: new Date('2025-11-27T14:05:30Z')
    //         },
    //         {
    //           origin: "recieved",
    //           type: "grade",
    //           value: "incorrect",
    //           sentTime: new Date('2025-11-27T14:05:30Z')
    //         },
    //         {
    //           origin: "sent",
    //           type: "message",
    //           value: "prosiiim",
    //           sentTime: new Date('2025-11-27T14:05:30Z')
    //         },
    //         {
    //           origin: "recieved",
    //           type: "grade",
    //           value: "correct",
    //           sentTime: new Date('2025-11-27T14:05:30Z')
    //         },
    //       ]
    //     },
    //     "ijfs": {
    //       id: "ijfs",
    //       name: "tohle je jmeno",
    //       text: "tohle jetextulohy",
    //       answer: "",
    //       diff: "A",
    //       images: [],
    //       focusedBy: [],
    //       chat: []
    //     },
    //     "alsAL": {
    //       id: "alsAL",
    //       name: "tohle je jmeno",
    //       text: "tohle jetextulohy",
    //       answer: "",
    //       diff: "A",
    //       images: [],
    //       focusedBy: ["adf", "ads"],
    //       chat: []
    //     },
    //     "PIJFSLKHBKJ": {
    //       id: "PIJFSLKHBKJ",
    //       name: "tohle je jmeno",
    //       text: "tohle jetextulohy",
    //       answer: "",
    //       diff: "A",
    //       images: [],
    //       focusedBy: [],
    //       chat: []
    //     },
    //
    //   }
    // });
    // currentState.update(state => {
    //   return {
    //     teamName: "Bambuláci 4. trida",
    //     money: 80,
    //     myId: "bablbamId",
    //     probsRemaining: [10, 2, -1],
    //     pricesBuy: [10, 30, 80],
    //     pricesSell: [10, 15, 40],
    //     procesSolve: [15, 50, 200]
    //   }
    // })
  });

  function handleBuy(diff: string) {
    socket.send(JSON.stringify({
      "name": "buy",
      "diff": diff,
    }));
  }
  function handleChat(probId: string, message: Omit<MessageType, 'sentTime'>) {
    console.log("Message of type '" + message.type + "': " + message.value);
  }
  function handleSell(probId: string) {
    socket.send(JSON.stringify({
      "name": "buy",
      "probid": probId,
    }));
  }
  function handleFocus(probId: string) {
    console.log("Focusing:", probId);
    socket.send(JSON.stringify({
      "name": "buy",
      "probid": probId,
    }));
  }
  function handleSolve(probId: string, answer: string) {
    socket.send(JSON.stringify({
      "name": "solve",
      "probid": probId,
      "answer": answer,
    }));
  }
</script>

<main>
  {#if $wsConnected}
    {#if $currentState.runningState === "before"}
      <Waitroom />
    {:else if $currentState.runningState === "running"}
      <Main buy={handleBuy} chat={handleChat} sell={handleSell} focus={handleFocus} solve={handleSolve} />
    {:else}
      <h1>GG, skončili jste {$currentState.rank}</h1>
      <!-- <Main buy={handleBuy} chat={handleChat} sell={handleSell} focus={handleFocus} solve={handleSolve} /> -->
      <!-- <Waitroom /> -->
      <!-- <h1>Soutěž skončila. Děkujeme za účast!</h1> -->
    {/if}
  {:else}
    <Disconnected/>
  {/if}
</main>

<style lang="scss">
  main {
    min-height: 0px;
    display: flex;
    flex-direction: column;
    width: 100%;
    height: 100%;
  }
</style>
