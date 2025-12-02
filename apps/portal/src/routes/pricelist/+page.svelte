<script lang="ts">
  // import type { PageData } from './$types';
  // export let data: PageData;
  let { data } = $props();
  const problems = Object.keys(data.prices.Buy);
  const categories = Object.keys(data.prices);
</script>

<main>
  <h1 class="title">Ceník</h1>

  <div class="overflow-x-auto">
    <table class="min-w-full bg-white border border-gray-200">
      <thead>
        <tr class="">
          <th class="py-2 px-4 border-b">Úloha</th>
          <th class="py-2 px-4 border-b">Koupit</th>
          <th class="py-2 px-4 border-b">Prodat</th>
          <th class="py-2 px-4 border-b">Vyřešit</th>
        </tr>
      </thead>
      <tbody>
        {#each problems as problem}
          <tr class="hover:bg-gray-50">
            <td class="py-2 px-4 border-b text-center">{problem}</td>
            {#each categories as category}
              <td class="py-2 px-4 border-b text-center">{data.prices[category][problem]}</td>
            {/each}
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</main>

<style lang="scss">
  main {
    display: flex;
    width: 100vw;
    height: 100vh;

    flex-direction: column;
    justify-content: center;
    align-items: center;

    .title {
      font-family: 'Lexend';
      font-size: 40px;
    }
    table {
      thead {
        th {
          font-family: 'Lexend';
          font-weight: 600;
          font-size: 25px;
          padding: 5px 10px;
        }
      }

      tbody {
        td {
          font-family: 'Lexend';
          font-weight: 400;
          font-size: 20px;
          padding: 5px 10px;
        }
      }

    }

  }

</style>