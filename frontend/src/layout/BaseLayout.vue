<template>
  <div class="flex h-full">
    <div class="h-full flex flex-col w-80 resize-x overflow-auto">
      <div class="h-[92px] flex-shrink-0 px-2 pt-[50px]" style="--wails-draggable: drag">
        <slot name="side-head"></slot>
      </div>
      <slot name="side-body"></slot>
    </div>
    <div class="h-full bg-systemBackground-light flex-grow flex flex-col justify-between pb-1">
      <slot name="main"></slot>
    </div>
  </div>
</template>
