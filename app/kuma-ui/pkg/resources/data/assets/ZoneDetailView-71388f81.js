import{d as k,u as z,q as l,s as p,o as a,a as s,w as f,h as c,b as _,g as h,k as b,e as y}from"./index-9d631905.js";import{_ as $}from"./ZoneDetails.vue_vue_type_script_setup_true_lang-45ff8dae.js";import{j as g,f as x,k as B,g as E,_ as V}from"./RouteView.vue_vue_type_script_setup_true_lang-76145142.js";import{_ as N}from"./RouteTitle.vue_vue_type_script_setup_true_lang-f639963c.js";import{_ as A}from"./EmptyBlock.vue_vue_type_script_setup_true_lang-255e2244.js";import{E as C}from"./ErrorBlock-be40f398.js";import{_ as D}from"./LoadingBlock.vue_vue_type_script_setup_true_lang-7f9cc3f9.js";import"./kongponents.es-bba90403.js";import"./CodeBlock.vue_vue_type_style_index_0_lang-9125ad7e.js";import"./DefinitionListItem-ad3ab377.js";import"./SubscriptionHeader.vue_vue_type_script_setup_true_lang-9b865501.js";import"./TabsWidget-0e0dd5da.js";import"./QueryParameter-70743f73.js";import"./TextWithCopyButton-6bd93ee0.js";import"./WarningsWidget.vue_vue_type_script_setup_true_lang-ffa4d4c0.js";const O={class:"zone-details"},T={key:3,class:"kcard-border","data-testid":"detail-view-details"},W=k({__name:"ZoneDetailView",setup(Z){const d=g(),e=z(),v=x(),{t:m}=B(),t=l(null),n=l(!0),o=l(null);p(()=>e.params.mesh,function(){e.name==="zone-cp-detail-view"&&i()}),p(()=>e.params.name,function(){e.name==="zone-cp-detail-view"&&i()}),w();function w(){v.dispatch("updatePageTitle",e.params.zone),i()}async function i(){n.value=!0,o.value=null;const u=e.params.zone;try{t.value=await d.getZoneOverview({name:u})}catch(r){t.value=null,r instanceof Error?o.value=r:console.error(r)}finally{n.value=!1}}return(u,r)=>(a(),s(V,null,{default:f(()=>[c(N,{title:_(m)("zone-cps.routes.item.title")},null,8,["title"]),h(),c(E,{breadcrumbs:[{to:{name:"zone-cp-list-view"},text:_(m)("zone-cps.routes.item.breadcrumbs")}]},{default:f(()=>[b("div",O,[n.value?(a(),s(D,{key:0})):o.value!==null?(a(),s(C,{key:1,error:o.value},null,8,["error"])):t.value===null?(a(),s(A,{key:2})):(a(),y("div",T,[c($,{"zone-overview":t.value},null,8,["zone-overview"])]))])]),_:1},8,["breadcrumbs"])]),_:1}))}});export{W as default};
