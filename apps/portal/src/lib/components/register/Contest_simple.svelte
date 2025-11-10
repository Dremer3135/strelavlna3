<script lang="ts">
  import type { Contest, Team } from "$lib/types/register";
  import { pocketbase } from "$lib/pocketbase";

  import block_orange from "$lib/assets/images/register/block_orange.svg?url";
  import block_yellow from "$lib/assets/images/register/block_yellow.svg?url";
  import block_purple from "$lib/assets/images/register/block_purple.svg?url";
  import block_pink from "$lib/assets/images/register/block_pink.svg?url";

  import strela from "$lib/assets/images/register/strela.svg?url";
  import vlna from "$lib/assets/images/register/vlna.svg?url";
  import { writable, type Writable } from "svelte/store";
  import { teamsStore, fetchTeams } from "$lib/stores/register";
  import { currentUser } from '$lib/stores/auth';
  import AddTeamModal from "./AddTeamModal.svelte";


  let { contest }: { contest: Contest } = $props();

  let showAddTeamModal = $state(false);

  function openAddTeamModal() {
    if ($currentUser && pocketbase.authStore.isValid && $currentUser?.verified) {
      showAddTeamModal = true;
    }
  }

  async function handleTeamCreated() {
    if ($currentUser && pocketbase.authStore.isValid && $currentUser?.verified) {
      await fetchTeams($currentUser.id);
    }
  }

  const filteredTeams = $derived($teamsStore.filter(team => team.contest === contest.id));

  const dates_raw_str = $derived([contest.onlineStart, contest.onSiteStart]);
  let dates_str: string[] = ["", ""];  // online, onSite
  let remaining_str = writable(["", ""]);

  $effect(() => {
    for (let i = 0; i < dates_raw_str.length; i++) {
      let date_american_numstr: string[] = new Date(dates_raw_str[i]).toLocaleDateString().split("/");
      let date_str = date_american_numstr[1] + ". " + date_american_numstr[0] + ". " + date_american_numstr[2];
      
      dates_str[i] = date_str;
    }
  });

  import { onMount } from "svelte";
    import RegisteredTeam from "./RegisteredTeam.svelte";

  onMount(() => {
    const interval = setInterval(() => {
      let remaining_new_str: string[] = ["", ""];
      for (let i = 0; i < dates_raw_str.length; i++) {
        let remaining_ms: number = new Date(dates_raw_str[i]).getTime() - Date.now();
        let remaining_sec: number, remaining_min: number, remaining_hr: number, remaining_day: number;
        
        if (remaining_ms > 0) {
          remaining_sec = Math.floor(remaining_ms / 1000 % 60);
          remaining_min = Math.floor(remaining_ms / 1000 / 60 % 60);
          remaining_hr = Math.floor(remaining_ms / 1000 / 60 / 60 % 24);
          remaining_day = Math.floor(remaining_ms / 1000 / 60 / 60 / 24);
        } else {
          remaining_sec = 0;
          remaining_min = 0;
          remaining_hr = 0;
          remaining_day = 0;
        }

        let remaining_sec_str: string = remaining_sec < 10 ? "0" + remaining_sec.toString() : remaining_sec.toString();
        let remaining_min_str: string = remaining_min < 10 ? "0" + remaining_min.toString() : remaining_min.toString();
        let remaining_hr_str: string = remaining_hr < 10 ? "0" + remaining_hr.toString() : remaining_hr.toString();
        let remaining_day_str: string = remaining_day.toString();

        remaining_new_str[i] = remaining_day_str + " dní " + remaining_hr_str + ":" + remaining_min_str + ":" + remaining_sec_str;
      }

      remaining_str.set(remaining_new_str);
    }, 1000);

    return () => clearInterval(interval);
  });
</script>

<main>
  {#if showAddTeamModal}
    <AddTeamModal 
      contest={contest} 
      on:close={() => showAddTeamModal = false} 
      on:created={handleTeamCreated} 
    />
  {/if}
  <div class="contest-container" on:click={openAddTeamModal}>
    <div class="tooltip">
      <h3>
        {#if $currentUser && pocketbase.authStore.isValid && $currentUser?.verified}
          Zaregistrovat tým
        {:else}
          Pro registraci týmu se přihlaste
        {/if}
      </h3>
    </div>
    <img src={block_orange} alt="" class="block" id="block-orange" />
    <img src={block_yellow} alt="" class="block" id="block-yellow" />
    <img src={block_purple} alt="" class="block" id="block-purple" />
    <img src={block_pink} alt="" class="block" id="block-pink" />
    <div class="contest-name">{contest.name}</div>
    <div class="contest-type">{contest.type == "math" ? "Matematika" : "Fyzika"}</div>
    <div class="date-holder">
      <p class="date">Online: {dates_str[0]}</p>
      <p class="remaining">{$remaining_str[0]}</p>
    </div>
    <div class="date-holder">
      <p class="date">Prezenční: {dates_str[1]}</p>
      <p class="remaining">{$remaining_str[1]}</p>
    </div>
    <img src={contest.type == "math" ? strela : vlna} class="contest-image">
  </div>
  
  <div class="registered-team-container">
    {#each filteredTeams as team}
      <RegisteredTeam team={team} />
    {/each}
  </div>
</main>

<style>
  main {
    width: 100%;
    display: flex;
    flex-direction: column;
    color: black;
  }

  .tooltip {
    position: absolute;
    top: -40px; /* Final position */
    left: 50%;
    opacity: 0;
    transform: translate(-50%, 20px); /* Start position (offset down) */
    background-color: #9500EB;
    padding: 5px 15px;
    border-radius: 5px;
    pointer-events: none;
    transition: transform 0.3s cubic-bezier(0.33, 1, 0.68, 1), opacity 0.3s cubic-bezier(0.33, 1, 0.68, 1);
    will-change: transform, opacity;
  }
  .tooltip h3 {
    margin: 0px;
    font-weight: 500;
    color: white;
  }

  .contest-container:hover .tooltip {
    opacity: 1;
    transform: translate(-50%, -10px); /* End position */
    animation: tooltip-animation 4s infinite steps(1, end);
  }

  .registered-team-container {
    display: flex;
    flex-direction: row;
    padding-inline: 100px;
    gap: 10px;
    flex-wrap: wrap;
  }

  .contest-container {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background-color: #f9f9f9;
    margin-bottom: 1rem;
    padding: 20px 30px;
    width: 100%;
    box-sizing: border-box;
    position: relative;
    color: #002c5e;
    min-height: 100px;
    outline-offset: 0px;
    outline: 5px dashed transparent;
    cursor: pointer;
  }

  .contest-container:hover {
    outline-color: #9500EB;
    animation: outline-animation 4s steps(1, end) infinite;

  }

  .block {
    position: absolute;
    animation: move_around 4s steps(1, end) infinite;
  }

  #block-orange {
    animation-delay: 0s;
  }

  #block-yellow {
    animation-delay: -1s;
  }

  #block-purple {
    animation-delay: -2s;
  }

  #block-pink {
    animation-delay: -3s;
  }

  @keyframes move_around {
    0% {
      top: -20px;
      left: -20px;
    }
    25% {
      top: -20px;
      left: calc(100% + 10px);
    }
    50% {
      top: calc(100% + 10px);
      left: calc(100% + 10px);
    }
    75% {
      top: calc(100% + 10px);
      left: -20px;
    }
    100% {
      top: -20px;
      left: -20px;
    }
  }

  .date-holder {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }
  .date-holder p {
    margin: 0px;
  }

  .date {
    font-size: 18px;
    font-weight: 600;
  }
  .remaining {
    font-size: 14;
    font-weight: 400;
    color: #6a83a0;
  }

  .contest-name {
    font-size: 20px;
    font-weight: 600;
  }

  @keyframes outline-animation {
    0% {
      outline-color: #9500EB;
    }
    25% {
      outline-color: #EB0072;
    }
    50% {
      outline-color: #EB6E00;
    }
    75% {
      outline-color: #EBAD00;
    }
  }

  @keyframes tooltip-animation {
    0% {
      background-color: #9500EB;
    }
    25% {
      background-color: #EB0072;
    }
    50% {
      background-color: #EB6E00;
    }
    75% {
      background-color: #EBAD00;
    }
  }

</style>
