<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import Main from './Main.svelte';
  import { probs } from '$lib/stores/probs';
  import { currentState, wsConnected } from '$lib/stores/state';
  import type { MessageType, ResultsAtom } from '$lib/types';
  import Disconnected from '$lib/assets/components/Disconnected.svelte';
  import Navbar from '$lib/assets/components/Navbar.svelte';

  let socket: WebSocket;
  
  onMount(() => {

    socket = new WebSocket(`https://sv.skrat.org/ws/corr?token=${$page.data.token}`);

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
          teamName: "",
          money: 0,
          myId: "",
          probsRemaining: [],
          pricesBuy: [],
          pricesSell: [],
          procesSolve: [],
          start: new Date(Date.now() + data.data.start),
          end: new Date(Date.now() + data.data.end),
          runningState: data.data.state,
          rank: -1,
          isAdmin: data.data.admin,
          results: undefined,
        });

        let nprobs: any = {};

        for (let tickid of Object.keys(data.data.bought_tickets)) {
          let tick = data.data.bought_tickets[tickid];
          let prob = tick.prob;
          nprobs[tick.team_id + ":" + prob.id] = {
            id: tick.team_id + ":" + prob.id,
            name: tick.team_name + ": " + prob.name,
            text: prob.text,
            answer: prob.answer,
            diff: prob.diff,
            images: prob.images,
            focusedBy: [],
            chat: data.data.tlines[tickid].map((e: any) => {return {
              origin: e.mside === "admin" ? "sent" : "recieved",
              type: e.mtype,
              value: e.msg,
              sentTime: new Date(e.time),
            }}),
            owned: "bought",
          };
        }

        for (let tickid of Object.keys(data.data.solved_tickets)) {
          let tick = data.data.solved_tickets[tickid];
          let prob = tick.prob;
          nprobs[tick.team_id + ":" + prob.id] = {
            id: tick.team_id + ":" + prob.id,
            name: prob.team_name,
            text: prob.text,
            answer: prob.answer,
            diff: prob.diff,
            images: prob.images,
            focusedBy: [],
            chat: data.data.tlines[tickid].map((e: any) => {return {
              origin: e.mside === "admin" ? "sent" : "recieved",
              type: e.mtype,
              value: e.msg,
              sentTime: new Date(e.time),
            }}),
            owned: "solved",
          };
        }

        for (let tickid of Object.keys(data.data.sold_tickets)) {
          let tick = data.data.sold_tickets[tickid];
          let prob = tick.prob;
          nprobs[tick.team_id + ":" + prob.id] = {
            id: tick.team_id + ":" + prob.id,
            name: prob.team_name,
            text: prob.text,
            answer: prob.answer,
            diff: prob.diff,
            images: prob.images,
            focusedBy: [],
            chat: data.data.tlines[tickid].map((e: any) => {return {
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
          currentProbs = Object.fromEntries(Object.entries(currentProbs).map(([_, prob]) => [prob.id, {...prob, focusedBy: prob.focusedBy.filter((x: any) => x !== $currentState.myId)}]));  // remove focus of specific plaeyer from all probs
          currentProbs[data.id].focusedBy.push($currentState.myId);  // add focus to specific prob
          return currentProbs;
        });
      } else if (data.name === "written") { 
        probs.update((e) => {
          Object.keys(e).forEach((k) => {
            if (k === data.teamid + ":" + data.probid) {
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
          e[data.teamid + ":" + data.prob.id] = {
            id: data.teamid + ":" + data.prob.id,
            name: data.teamname + ": " + data.prob.name,
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
        // currentState.update((e) => { return {...e, money: data.money}})
      } else if (data.name === "sold") {
        probs.update((e) => {
          e[data.teamid + ":" + data.probid].owned = "sold";
          // currentState.update((e) => { return {...e, money: data.money}})
          return e;
        });
      } else if (data.name === "solved") {
        probs.update((e) => {
          e[data.teamid + ":" + data.probid].owned = "solved";
          // currentState.update((e) => { return {...e, money: data.money}})
          return e;
        });
      } else if (data.name === "start") {
        currentState.update((e) => { return {...e, runningState: "running"}})
      } else if (data.name === "results") {
        let nresults: Record<string, ResultsAtom> = {};
        for (let k of Object.keys(data.data)) {
          let rec = data.data[k];
          nresults[k] = {
            teamName: rec.name,
            rank: rec.rank,
            money: rec.money,
          }
        }
        currentState.update((e) => { return {...e, runningState: "results", results: nresults}})
      }
    });

  });

  function handleChat(probId: string, message: Omit<MessageType, 'sentTime'>) {
    if (message.type === "grade") {
      socket.send(JSON.stringify({
        "name": "grade",
        "id": probId,
        "decision": message.value,
      }));
      return;
    }

    socket.send(JSON.stringify({
      "name": "write",
      "id": probId,
      "message": message.value,
      "mtype": message.type,
    }));

    // console.log("Message of type '" + message.type + "': " + message.value);
  }
  function handleFocus(probId: string) {
    socket.send(JSON.stringify({"name": "focus", "id": probId}));
    // console.log("Focusing:", probId);
  }
  function handleGrade(probId: string) {
    socket.send(JSON.stringify({"name": "grade", "id": probId}));
    console.log("Grading:", probId);
  }
  function handleStart() {
    console.log("start");
    socket.send(JSON.stringify({"name": "start"}));
  }
  function handleEnd() {
    socket.send(JSON.stringify({"name": "end"}));
  }
  function handleResults() {
    socket.send(JSON.stringify({"name": "results"}));
  }
  function handlePause() {
    socket.send(JSON.stringify({"name": "pause"}));
    console.log("pause");
  }
  function handleResume() {
    socket.send(JSON.stringify({"name": "resume"}));
    console.log("resume");
  }

  $effect(() => {
    console.log("lalalala:O", Object.values($probs));
  });
</script>


{#if $wsConnected}
  <Navbar start={handleStart} end={handleEnd} results={handleResults} pause={handlePause} resume={handleResume}/>
{/if}
<main>
  {#if $wsConnected}
    {#if $currentState.runningState == "results"}
      <ul>
        {#each Object.values($currentState.results ?? {}).toSorted((a, b) => a.rank - b.rank) as item}
          <li>{item.rank}: {item.teamName} ({item.money})</li>
        {/each}
      </ul>
    {:else}
      <Main chat={handleChat} focus={handleFocus} />
    {/if}
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
