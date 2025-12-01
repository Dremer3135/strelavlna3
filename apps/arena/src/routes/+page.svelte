<script lang="ts">
  import { onMount } from 'svelte';
  import Main from './Main.svelte';
  import { probs } from '$lib/stores/probs';
  import { currentState, wsConnected } from '$lib/stores/state';
  import type { MessageType, CurrentState } from '$lib/types';
  import Waitroom from './Waitroom.svelte';
  import Disconnected from '$lib/assets/components/Disconnected.svelte';
  import Navbar from '$lib/assets/components/Navbar.svelte';
  import Paused from '$lib/assets/components/Paused.svelte';

  let { data } = $props();

  let socket: WebSocket;
  
  onMount(() => {
    socket = new WebSocket(`https://sv.skrat.org/ws?token=${data.token}`);

    socket.onopen = () => {
      wsConnected.set(true);
    };
    socket.onclose = () => {
      wsConnected.set(false);
    }
    socket.onerror = () => {
      wsConnected.set(false);
    }

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
              sentTime: new Date(e.time),
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
              sentTime: new Date(e.time),
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
              sentTime: new Date(e.time),
            }}),
            owned: "sold",
          };
        }

        probs.set(nprobs);
      } else if (data.name === "focus") {
        probs.update((currentProbs) => {
          currentProbs = Object.fromEntries(Object.entries(currentProbs).map(([_, prob]) => [prob.id, {...prob, focusedBy: prob.focusedBy.filter((x: any) => x !== data.playerid)}]));  // remove focus of specific plaeyer from all probs
          currentProbs[data.probid].focusedBy.push(data.playerid);  // add focus to specific prob
          return currentProbs;
        });
      } else if (data.name === "written") {
        probs.update((e) => {
          Object.keys(e).forEach((k) => {
            if (k === data.probid) {
              e[k].chat.push({
                origin: data.origin,
                type: data.type,
                value: data.text,
                sentTime: new Date(data.time),
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
          currentState.update((e) => { return {...e, money: data.money}});
          return e;
        });
      } else if (data.name === "solved") {
        probs.update((e) => {
          e[data.probid].owned = "solved";
          currentState.update((e) => { return {...e, money: data.money}})
          return e;
        });
      } else if (data.name === "start") {
        currentState.update((e) => { return {...e, runningState: "running"}})
      } else if (data.name === "end") {
        currentState.update((e) => { return {...e, runningState: "after"}})
      } else if (data.name === "results") {
        currentState.update((e) => { return {...e, runningState: "results", money: data.money, rank: data.rank}})
      } else if (data.name === "paused") {
        currentState.update((e) => { return {...e, runningState: "paused"}})
      } else if (data.name === "resumed") {
        currentState.update((e) => { return {...e, runningState: "running"}})
      } else if (data.name === "chmoney") {
        currentState.update((e) => { return {...e, money: data.money}})
      }



      // console.log(data);
      // console.log($probs);
      // console.log($currentState);
    });
  });

  function handleBuy(diff: string) {
    console.log(`Buying ${diff}`);
    socket.send(JSON.stringify({
      "name": "buy",
      "diff": diff,
    }));
  }
  function handleChat(probId: string, message: Omit<MessageType, 'sentTime'>) {
    console.log("Message of type '" + message.type + "': " + message.value);
    console.log(`Writing ${probId} ${message}`);
    socket.send(JSON.stringify({
      "name": "write",
      "probid": probId,
      "message": message.value,
      "mtype": message.type,
    }));
  }
  function handleSell(probId: string) {
    console.log(`Selling ${probId}`);
    socket.send(JSON.stringify({
      "name": "sell",
      "probid": probId,
    }));
  }
  function handleFocus(probId: string) {
    console.log(`Focusing ${probId}`);
    console.log("Focusing:", probId);
    socket.send(JSON.stringify({
      "name": "focus",
      "probid": probId,
    }));
  }
  function handleSolve(probId: string, answer: string) {
    console.log(`Solving ${probId} ${answer}`);
    socket.send(JSON.stringify({
      "name": "solve",
      "probid": probId,
      "answer": answer,
    }));
  }
</script>

{#if $wsConnected}
  <Navbar />
{/if}
<main>
  {#if $wsConnected}
    {#if $currentState.runningState === "before"}
      <Waitroom />
    {:else if $currentState.runningState === "running"}
      <Main buy={handleBuy} chat={handleChat} sell={handleSell} focus={handleFocus} solve={handleSolve} />
    {:else if $currentState.runningState === "after"}
      <h1>Čekáme na výsledky...</h1>
    {:else if $currentState.runningState === "results"}
      <h1>GG, skončili jste {$currentState.rank}. s {$currentState.money} body</h1>
      <h2>Pokud jste v top 15, {$currentState.rank > 15 ? "což nejste, " : ""}očekávejte email o informacích o prezenčním kole</h2> 
      <h2>Úlohy připravila třída 7.M</h2>
      <h2 style="color: blue; font-family:Impact, Haettenschweiler, 'Arial Narrow Bold', sans-serif">Zážitek ze hry vám přináší tým Skrat</h2>
      <h2 style="color: green">Veškeré technické potíže způsobilo ChatGPT a Gemini...</h2>
      <h2 style="color: red">Na závěrečnou obrazovku nezbyl budget (ne že by nějaký byl)</h2>
      <h2>Užijte si alespoň tyto obrázky koťátek</h2>
      {#each Array(7).fill(0) as _item, i}
        <img src="https://cataas.com/cat?t={i}" alt="kočka">
      {/each}
      <!-- <Main buy={handleBuy} chat={handleChat} sell={handleSell} focus={handleFocus} solve={handleSolve} /> -->
      <!-- <Waitroom /> -->
      <!-- <h1>Soutěž skončila. Děkujeme za účast!</h1> -->
    {:else}
      <Paused isPaused={true} />
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
