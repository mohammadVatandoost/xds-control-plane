import{d as h,u as _,c as f,o as i,a as m,w as p,h as r,b as c,g as y,f as d}from"./index-9d631905.js";import{_ as P}from"./PolicyDetails.vue_vue_type_script_setup_true_lang-5a7aa06f.js";import{f as b,k as x,g as N,_ as k}from"./RouteView.vue_vue_type_script_setup_true_lang-76145142.js";import{_ as w}from"./RouteTitle.vue_vue_type_script_setup_true_lang-f639963c.js";import"./StatusInfo.vue_vue_type_script_setup_true_lang-ea244d88.js";import"./EmptyBlock.vue_vue_type_script_setup_true_lang-255e2244.js";import"./kongponents.es-bba90403.js";import"./ErrorBlock-be40f398.js";import"./LoadingBlock.vue_vue_type_script_setup_true_lang-7f9cc3f9.js";import"./ResourceCodeBlock.vue_vue_type_script_setup_true_lang-5d930ce7.js";import"./CodeBlock.vue_vue_type_style_index_0_lang-9125ad7e.js";import"./TextWithCopyButton-6bd93ee0.js";import"./toYaml-4e00099e.js";import"./TabsWidget-0e0dd5da.js";import"./QueryParameter-70743f73.js";const G=h({__name:"PolicyDetailView",props:{mesh:{},policyPath:{},policyName:{}},setup(l){const e=l,n=_(),t=b(),{t:a}=x(),o=f(()=>t.state.policyTypesByPath[e.policyPath]);u();function u(){t.dispatch("updatePageTitle",n.params.policy)}return(T,V)=>(i(),m(k,{module:"policies"},{default:p(({route:s})=>[r(w,{title:c(a)("policies.routes.item.title")},null,8,["title"]),y(),r(N,{breadcrumbs:[{to:{name:"policies-list-view",params:{mesh:s.params.mesh,policyPath:s.params.policyPath}},text:c(a)("policies.routes.item.breadcrumbs")}]},{default:p(()=>[o.value?(i(),m(P,{key:0,name:e.policyName,mesh:e.mesh,path:e.policyPath,type:o.value.name},null,8,["name","mesh","path","type"])):d("",!0)]),_:2},1032,["breadcrumbs"])]),_:1}))}});export{G as default};