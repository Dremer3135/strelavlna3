<script lang="ts">
  import { onMount } from 'svelte';
  import { wsConnected } from '$lib/stores/state';
  import { adminViewedTeams } from '$lib/stores/adminView';

  let { data } = $props();

  let socket: WebSocket;
  
  onMount(() => {

    socket = new WebSocket(`https://sv.skrat.org/ws/admin?token=${data.token}`);

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
        adminViewedTeams.set(data.data);

      } else if (data.name === "tchange") {
        adminViewedTeams.update((currentAVT) => {
          currentAVT[data.teamid] = { ...(currentAVT[data.teamid] ?? {}), ...data };
          return currentAVT;
        });
      }
    });
  });

  let sortedTeams = $derived(Object.entries($adminViewedTeams).sort(([, a], [, b]) => b.money - a.money));

</script>

<main>
  {#each sortedTeams as [id, team]}
    {#key id}
      <button class="team" onclick={() => {
        navigator.clipboard.writeText(team.teamid);
      }}>
        <!-- <h3>{team.name}</h3> -->
        <p>{team.money}</p>
        <p class="id">{team.teamname}</p>
      </button>
    {/key}
  {/each}
</main>

<style lang="scss">
  main {
    padding: 100px;
    display: flex;
    flex-wrap: wrap;
    gap: 15px;

    .team {
      all: unset;
      cursor: pointer;
      overflow: visible;
      position: relative;
      display: flex;
      border: var(--color-yellow) 3px solid;
      background-color: color-mix(in srgb, var(--color-yellow) 5%, transparent 95%);
      padding: 10px 20px;
      box-sizing: border-box;
      border-radius: 5px;
      animation: boom 1s cubic-bezier(0.215, 0.610, 0.355, 1) forwards;

      h3 {
        font-family: 'Lexend';
        font-weight: 600;
        font-size: 25px;
      }

      p {
        font-family: 'Lexend';
        font-weight: 500;
        font-size: 18px;
        margin: 0px;
        display: block;
      }

      .id {
        position: absolute;
        bottom: 0px;
        left: 50%;
        transform: translate(-50%, 0%);
        opacity: 0;
        transition: all cubic-bezier(0.215, 0.610, 0.355, 1) 0.3s;
        background-color: var(--color-yellow);
        box-shadow: 0px 0px 0px 5px white;
        padding: 5px 10px;
        border-radius: 5px;
        box-sizing: border-box;
        text-wrap: nowrap;
        font-size: 11px;
        pointer-events: none;
      }

      &:hover{
        .id {
          opacity: 1;
          transform: translate(-50%, calc(100% + 15px));
        }
      }
    }
  }

  @keyframes boom {
    0% {
      border-color: var(--color-orange);
      background-color: var(--color-orange);
      color: white;
    }
    100% {
      border-color: var(--color-yellow);
      background-color: color-mix(in srgb, var(--color-yellow) 5%, transparent 95%);
      color: black;
    }
  }
</style>