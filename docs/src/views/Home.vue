<template>
    <div class="bg-white text-slate-700 dark:bg-slate-900 dark:text-slate-400">
        <LanguageSwitch />
        <ThemeSwitch />
        <div class="flex justify-center items-center min-h-screen flex-col">
            <div class="my-20">
                <div class="flex mx-auto max-w-fit">
                    <h1 class="text-xl flex-1 sm:text-4xl font-extrabold text-slate-900 tracking-tight
                        dark:text-slate-200 text-right whitespace-nowrap">
                        {{ $t('h1') }}
                    </h1>
                    <div class="flex-none sm:mt-2 mt-1 ml-2">
                        <LabeledSwitch ref="versionSwitch" />
                    </div>
                </div>
                <div
                    class="mx-auto sm:max-w-2xl max-w-[320px] text-center pb-6 pt-1 text-lg text-slate-700 dark:text-slate-400">
                    <i18n-t keypath="desc" scope="global">
                        <template v-slot:input>
                            <span
                                class="text-sm tracking-tight dark:text-slate-300 text-slate-500">[PNG,GIF,JPG,JPEG]</span>
                        </template>
                    </i18n-t>
                </div>
                <Urls :data="urls" :cli="versionSwitch.checked" />
            </div>
        </div>
    </div>
</template>
<script setup lang="ts">
import { useHead } from '@unhead/vue';
import { ref } from 'vue';
import Urls from '@/components/Urls.vue';
import LabeledSwitch from '@/components/LabeledSwitch.vue'
import LanguageSwitch from '@/components/LanguageSwitch.vue';
import ThemeSwitch from '@/components/ThemeSwitch.vue';
const versionSwitch = ref<{ checked: boolean }>({ checked: false })

const urls = ref([
    { title: "TO PNG", url: "toformat.link/png=", data: "http://toformat.link/img.jpg", curlEnd: " > img.png" },
    { title: "TO GIF", url: "toformat.link/gif=", data: "http://toformat.link/img.jpg", curlEnd: " > img.gif" },
    { title: "TO JPG", url: "toformat.link/jpg=", data: "http://toformat.link/img.png", curlEnd: " > img.jpg" },
    { title: "TO JPEG", url: "toformat.link/jpeg=", data: "http://toformat.link/img.png", curlEnd: " > img.jpeg" },
])
useHead({
    title: "toformat docs",
    meta: [
        {
            name: 'description',
            content: 'toformat documentation'
        },
        {
            name: 'robots',
            content: 'index'
        },
    ]
})
</script>
