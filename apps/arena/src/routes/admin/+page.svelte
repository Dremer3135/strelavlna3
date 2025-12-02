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

      } else if (data.name === "change") {
        adminViewedTeams.update((currentAVT) => {
          currentAVT[data.teamid] = { ...currentAVT[data.teamid], ...data.change };
          return currentAVT;
        });
      }
    });
  });

</script>

<main>
  {#each Object.entries($adminViewedTeams) as [id, team]}
    <div class="team">
      <h3>{team.name}</h3>
      <p>{team.money} DC</p>
    </div>
  {/each}
</main>

<style lang="scss">
  main {
    padding: 100px;
    display: flex;
    flex-wrap: wrap;

    .team {
      display: flex;
      border: var(--color-yellow) 3px solid;
      background-color: color-mix(in srgb, var(--color-yellow) 5%, transparent 95%);
      padding: 10px 20px;

      h3 {
        font-family: 'Lexend';
        font-weight: 600;
        font-size: 25px;
      }

      p {
        font-family: 'Lexend';
        font-weight: 500;
        font-size: 18px;
      }
    }
  }
</style>