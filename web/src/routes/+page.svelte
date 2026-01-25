<script lang="ts">
  import * as m from "$lib/paraglide/messages"
</script>

<div class="flex w-full h-full justify-start items-center flex-col gap-15">
  <div class="flex justify-center items-center flex-col relative top-10">
    <h1 class="text-5xl font-bold">Pdnode</h1>
    <p class="py-6 text-gray-500 text-center">
      {m.company_description()}
    </p>
  </div>

  <div class="flex justify-center items-center flex-col">
    <h3 class="text-3xl font-bold" >{m["about-us-title"]()}</h3>
    <p class="py-6 text-gray-500 text-center leading-loose whitespace-pre-line">
      {m.about_us()}
    </p>
  </div>

  <div class="flex justify-center items-center flex-col">
    <h3 class="text-3xl font-bold m-10">{m["team-member"]()}</h3>
    <div class="flex flex-col items-center gap-4 md:flex-row md:flex-wrap md:justify-center">
    <div class="card w-96 bg-base-100 card-md shadow-sm">
      <div class="card-body">
        <h2 class="card-title">Bret</h2>
        <p class="text-gray-600">{m["join-time"]()}: 2024-07-10</p>
        <p class="text-gray-600">{m.position()}: {m.founder()}</p>
        <p class="text-gray-600">{m["main-tasks"]()}: {m["main-tasks_full"]()}</p>
      </div>
    </div>
    <div class="card w-96 bg-base-100 card-md shadow-sm">
      <div class="card-body">
        <h2 class="card-title">Mark</h2>
        <p class="text-gray-600">{m["join-time"]()}: ????</p>
        <p class="text-gray-600">{m.position()}: {m.admin()}</p>
        <p class="text-gray-600">{m["main-tasks"]()}: {m["main-tasks_front"]()}</p>
      </div>
    </div>
    </div>
  </div>

  <div class="flex justify-center items-center flex-col gap-4 relative bottom-5">
    <h3 class="text-3xl font-bold">{m["contact-us"]()}</h3>
    <ul class="flex justify-center items-center flex-col gap-1">
      <li class="text-gray-600"><b class=" mr-1">{m.email()}: </b>bret@pdnode.com</li>
      <li class="link text-blue-500"><a href="https://github.com/BretRen/pdnode/issues">Github Issue</a></li>
    </ul>
  </div>
  
  
</div>