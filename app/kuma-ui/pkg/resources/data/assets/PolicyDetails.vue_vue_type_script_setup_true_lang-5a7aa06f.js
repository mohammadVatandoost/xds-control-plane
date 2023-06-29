import{d as k,q as p,c as P,s as b,G as D,r as x,o as y,a as N,w as t,g as u,R as T,S as B,k as w,e as v,j as V,h as n,t as g,F as A}from"./index-9d631905.js";import{_ as L}from"./StatusInfo.vue_vue_type_script_setup_true_lang-ea244d88.js";import{j as C}from"./RouteView.vue_vue_type_script_setup_true_lang-76145142.js";import{_ as $}from"./ResourceCodeBlock.vue_vue_type_script_setup_true_lang-5d930ce7.js";import{T as E}from"./TabsWidget-0e0dd5da.js";import{T as F}from"./TextWithCopyButton-6bd93ee0.js";const j=w("h2",null,"Dataplanes",-1),I=k({__name:"PolicyConnections",props:{mesh:{type:String,required:!0},policyPath:{type:String,required:!0},policyName:{type:String,required:!0}},setup(f){const e=f,_=C(),m=p(!1),o=p(!0),r=p(!1),i=p([]),a=p(""),l=P(()=>{const h=a.value.toLowerCase();return i.value.filter(({dataplane:s})=>s.name.toLowerCase().includes(h))});b(()=>e.policyName,function(){d()}),D(function(){d()});async function d(){r.value=!1,o.value=!0;try{const{items:h,total:s}=await _.getPolicyConnections({mesh:e.mesh,policyPath:e.policyPath,policyName:e.policyName});m.value=s>0,i.value=h??[]}catch{r.value=!0}finally{o.value=!1}}return(h,s)=>{const S=x("router-link");return y(),N(L,{"has-error":r.value,"is-loading":o.value,"is-empty":!m.value},{default:t(()=>[j,u(),T(w("input",{id:"dataplane-search","onUpdate:modelValue":s[0]||(s[0]=c=>a.value=c),type:"text",class:"k-input mt-4",placeholder:"Filter by name",required:"","data-testid":"dataplane-search-input"},null,512),[[B,a.value]]),u(),(y(!0),v(A,null,V(l.value,(c,q)=>(y(),v("p",{key:q,class:"mt-2","data-testid":"dataplane-name"},[n(S,{to:{name:"data-plane-detail-view",params:{mesh:c.dataplane.mesh,dataPlane:c.dataplane.name}}},{default:t(()=>[u(g(c.dataplane.name),1)]),_:2},1032,["to"])]))),128))]),_:1},8,["has-error","is-loading","is-empty"])}}}),M={class:"policy-details kcard-border"},R={class:"entity-heading","data-testid":"policy-single-entity"},z=k({__name:"PolicyDetails",props:{mesh:{type:String,required:!0},path:{type:String,required:!0},name:{type:String,required:!0},type:{type:String,required:!0}},setup(f){const e=f,_=C(),m=[{hash:"#overview",title:"Overview"},{hash:"#affected-dpps",title:"Affected DPPs"}],o=P(()=>({name:"policy-detail-view",params:{mesh:e.mesh,policy:e.name,policyPath:e.path}}));async function r(i){const{name:a,mesh:l,path:d}=e;return await _.getSinglePolicyEntity({name:a,mesh:l,path:d},i)}return(i,a)=>{const l=x("router-link");return y(),v("div",M,[n(E,{tabs:m},{tabHeader:t(()=>[w("h1",R,[u(g(e.type)+`:

          `,1),n(F,{text:e.name},{default:t(()=>[n(l,{to:o.value},{default:t(()=>[u(g(e.name),1)]),_:1},8,["to"])]),_:1},8,["text"])])]),overview:t(()=>[n($,{id:"code-block-policy","resource-fetcher":r,"resource-fetcher-watch-key":e.name,"is-searchable":""},null,8,["resource-fetcher-watch-key"])]),"affected-dpps":t(()=>[n(I,{mesh:e.mesh,"policy-name":e.name,"policy-path":e.path},null,8,["mesh","policy-name","policy-path"])]),_:1})])}}});export{z as _};
