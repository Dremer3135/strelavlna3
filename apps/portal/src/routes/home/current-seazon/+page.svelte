<script lang="ts">
  import { pocketbase } from "$lib/pocketbase";
  import { currentUser } from "$lib/stores/auth";
  import { contests } from "$lib/stores/contests";
  import ContestSimple from "./ContestSimple.svelte";
  import Auth from "./Auth.svelte";
  import TeamDetailsModal from "./TeamDetailsModal.svelte";
  import blocks_clumped from "$lib/assets/images/general/blocks_clumped.svg";
  import Navbar from "$lib/components/general/Navbar.svelte";
  import { myTeamsStore } from "$lib/stores/myTeams";
  // import TeamsPanel from "./TeamsPanel.svelte";

  let showAuth = false;
  let authType: "login" | "register" = "login";
  let showTeamsPanel = false;
  let showAddTeam = false;
  let showTeamDetails = false;
  let selectedTeam: any = null;

  let { data } = $props();

  contests.set(data.contests);
  myTeamsStore.set(data.teams)


  $effect(() => {
    console.log($contests);
  });


</script>

<div>
  
  {#if showAuth}
    <Auth on:close={() => (showAuth = false)} type={authType} />
  {/if}

  <!-- {#if showTeamsPanel}
    <TeamsPanel 
      on:close={() => (showTeamsPanel = false)} 
      on:showAddTeam={() => (showAddTeam = true)}
      on:showTeamDetails={(e) => {selectedTeam = e.detail; showTeamDetails = true;}}
    />
  {/if} -->
<!-- 
  {#if showAddTeam}
      <AddTeamModal on:close={() => (showAddTeam = false)} />
  {/if} -->

  {#if showTeamDetails}
      <TeamDetailsModal team={selectedTeam} on:close={() => (showTeamDetails = false)} />
  {/if}

  <div class="content-holder">
    <div class="logout-page">
      <h1 class="upcomming"><span class="anchor"><img src={blocks_clumped} class="blocks-clumped"></span> Nadcházející</h1>
      <div class="contests-list">
        {#each $contests as contest}
          <ContestSimple {contest} />
        {/each}
      </div>
    </div>
  </div>
</div>

<style>
  .contests-list {
    margin-top: 2rem;
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 100px;
  }

  .logout-page {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    padding: 80px;
    max-width: 1300px; 
    width: 100%;
    box-sizing: border-box;
  }

  .upcomming {
    font-family: "Fredoka";
    font-size: 50px;
    font-weight: bold;
    color: #002C5E;
    font-weight: 650;

    margin-bottom: 100px;
  }

  .content-holder {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }

  .anchor {
    position: relative;
  }

  .blocks-clumped {
    position: absolute;
    top: 55px;
    left: -20px;
  }
  @media (max-width: 600px) {
    .upcomming {
      font-size: 30px;
      margin-bottom: 50px;
      margin-left: 20px;
    }
    .logout-page {
      padding: 40px;
    }
  }
</style>
