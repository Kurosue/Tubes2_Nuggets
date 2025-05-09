"use client";
import React, { useCallback, useMemo, useRef, useState } from "react";
import { createPortal } from "react-dom";
import * as d3 from "d3";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { useIsomorphicLayoutEffect } from "@/lib/utils";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";

type D3CanvasBaseNodeDatum = d3.SimulationNodeDatum & { id: string };
type D3CanvasBaseLinkDatum<NodeDatum extends D3CanvasBaseNodeDatum> = d3.SimulationLinkDatum<NodeDatum> & { id: string, source: NodeDatum, target: NodeDatum };
type D3CanvasNodeRenderer<NodeDatum extends D3CanvasBaseNodeDatum, T extends Element> = 
	{
		enter: (selection: d3.Selection<d3.EnterElement, NodeDatum, SVGGElement, unknown>) => d3.Selection<T, NodeDatum, SVGGElement, unknown>,
		update: (selection: d3.Selection<T, NodeDatum, SVGGElement, unknown>) => void,
		exit: (selection: d3.Selection<T, NodeDatum, SVGGElement, unknown>) => void
	};
type D3CanvasLinkRenderer<LinkDatum extends D3CanvasBaseLinkDatum<any>, T extends Element> = 
	{
		enter: (selection: d3.Selection<d3.EnterElement, LinkDatum, SVGGElement, unknown>) => d3.Selection<T, LinkDatum, SVGGElement, unknown>,
		update: (selection: d3.Selection<T, LinkDatum, SVGGElement, unknown>) => void,
		exit: (selection: d3.Selection<T, LinkDatum, SVGGElement, unknown>) => void
	};
const newD3CanvasHandler = <
	NodeDatum extends D3CanvasBaseNodeDatum,
	LinkDatum extends D3CanvasBaseLinkDatum<NodeDatum>,
	NodeElement extends Element,
	LinkElement extends Element
>(
	svgElement: SVGSVGElement
) => {
	let nodeRenderer: D3CanvasNodeRenderer<NodeDatum, NodeElement> | null = null;
	let linkRenderer: D3CanvasLinkRenderer<LinkDatum, LinkElement> | null = null;
	const svg = d3.select(svgElement);
	const container = svg.append("g");
	const zoomer = d3.zoom<SVGSVGElement, unknown>().scaleExtent([0.1, 3])
		.on("zoom", e => container.attr("transform", e.transform));
	svg.call(zoomer);
	const linkGroup = container.append("g").attr("class", "links");
	const nodeGroup = container.append("g").attr("class", "nodes");
	const allNodes = new Map<string, NodeDatum>();
	const allLinks = new Map<string, LinkDatum>();
	let allNodesSelector: d3.Selection<NodeElement, NodeDatum, SVGGElement, unknown> | null = null;
	let allLinksSelector: d3.Selection<LinkElement, LinkDatum, SVGGElement, unknown> | null = null;
	const simulation = d3.forceSimulation<NodeDatum>()
		.force("charge", d3.forceManyBody<NodeDatum>().strength(-200))
		.force("link", d3.forceLink<NodeDatum, LinkDatum>().id(d => d.id).distance(100))
		.on("tick", () => {
			if(allNodesSelector != null && nodeRenderer != null)
				nodeRenderer.update(allNodesSelector);
			if(allLinksSelector != null && linkRenderer != null)
				linkRenderer.update(allLinksSelector);
		});
	const dragger = d3.drag<NodeElement, NodeDatum>()
		.on("start", (e, d) => {
			if(!e.active) simulation.alphaTarget(0.3).restart();
			d.fx = e.x;
			d.fy = e.y;
		})
		.on("drag", (e, d) => {
			d.fx = e.x;
			d.fy = e.y;
		})
		.on("end", (e, d) => {
			if(!e.active) simulation.alphaTarget(0);
			d.fx = null;
			d.fy = null;
		});
	const refreshData = () => {
		const nodes = nodeGroup.selectAll<NodeElement, NodeDatum>("& > *").data(allNodes.values(), d => d.id);
		const nodesEnter = nodeRenderer != null ? nodeRenderer.enter(nodes.enter()) : 
			d3.select(null) as any as d3.Selection<NodeElement, NodeDatum, SVGGElement, unknown>;
		nodesEnter.call(dragger);
		const nodesMerge = nodesEnter.merge(nodes);
		if(nodeRenderer != null)
			nodeRenderer.update(nodesMerge);
		const nodesExit = nodes.exit<NodeDatum>();
		if(nodeRenderer != null)
			nodeRenderer.exit(nodesExit);
		nodesExit.remove();

		const links = linkGroup.selectAll<LinkElement, LinkDatum>("& > *").data(allLinks.values(), d => d.id);
		const linksEnter = linkRenderer != null ? linkRenderer.enter(links.enter()) : 
			d3.select(null) as any as d3.Selection<LinkElement, LinkDatum, SVGGElement, unknown>;
		const linksMerge = linksEnter.merge(links);
		if(linkRenderer != null)
			linkRenderer.update(linksMerge);
		const linksExit = links.exit<LinkDatum>();
		if(linkRenderer != null)
			linkRenderer.exit(linksExit);
		linksExit.remove();

		simulation.nodes([...allNodes.values()]);
		(simulation.force("link") as d3.ForceLink<NodeDatum, LinkDatum>).links([...allLinks.values()]);
		simulation.alpha(0.6).restart();

		allNodesSelector = nodesMerge;
		allLinksSelector = linksMerge;
	};
	return {
		svgElement,
		get nodeRenderer() { return nodeRenderer },
		set nodeRenderer(v) {
			if(v == nodeRenderer) return;
			{
				const nodes = nodeGroup.selectAll<NodeElement, NodeDatum>("& > *").data([] as NodeDatum[], d => d.id);
				const nodesExit = nodes.exit<NodeDatum>();
				if(nodeRenderer != null)
					nodeRenderer.exit(nodesExit);
				nodesExit.remove();
			}
			nodeRenderer = v;
			{
				const nodes = nodeGroup.selectAll<NodeElement, NodeDatum>("& > *").data(allNodes.values(), d => d.id);
				const nodesEnter = nodeRenderer != null ? nodeRenderer.enter(nodes.enter()) : 
					d3.select(null) as any as d3.Selection<NodeElement, NodeDatum, SVGGElement, unknown>;
				nodesEnter.call(dragger);
				const nodesMerge = nodesEnter.merge(nodes);
				if(nodeRenderer != null)
					nodeRenderer.update(nodesMerge);
				const nodesExit = nodes.exit<NodeDatum>();
				if(nodeRenderer != null)
					nodeRenderer.exit(nodesExit);
				nodesExit.remove();
				allNodesSelector = nodesMerge;
			}
		},
		get linkRenderer() { return linkRenderer },
		set linkRenderer(v) {
			if(v == linkRenderer) return;
			{
				const links = linkGroup.selectAll<LinkElement, LinkDatum>("& > *").data([] as LinkDatum[], d => d.id);
				const linksExit = links.exit<LinkDatum>();
				if(linkRenderer != null)
					linkRenderer.exit(linksExit);
				linksExit.remove();
			}
			linkRenderer = v;
			{
				const links = linkGroup.selectAll<LinkElement, LinkDatum>("& > *").data(allLinks.values(), d => d.id);
				const linksEnter = linkRenderer != null ? linkRenderer.enter(links.enter()) : 
					d3.select(null) as any as d3.Selection<LinkElement, LinkDatum, SVGGElement, unknown>;
				const linksMerge = linksEnter.merge(links);
				if(linkRenderer != null)
					linkRenderer.update(linksMerge);
				const linksExit = links.exit<LinkDatum>();
				if(linkRenderer != null)
					linkRenderer.exit(linksExit);
				linksExit.remove();
				allLinksSelector = linksMerge;
			}
		},
		svg,
		container,
		zoomer,
		linkGroup,
		nodeGroup,
		allNodes,
		allLinks,
		get allNodesSelector() { return allNodesSelector },
		set allNodesSelector(v) { allNodesSelector = v },
		get allLinksSelector() { return allLinksSelector },
		set allLinksSelector(v) { allLinksSelector = v },
		simulation,
		dragger,
		refreshData
	};
};
type D3CanvasRefType<NodeDatum extends D3CanvasBaseNodeDatum, LinkDatum extends D3CanvasBaseLinkDatum<NodeDatum>, NodeElement extends Element, LinkElement extends Element> = 
	SVGSVGElement & { handler: ReturnType<typeof newD3CanvasHandler<NodeDatum, LinkDatum, NodeElement, LinkElement>> };
const D3Canvas = <
	NodeDatum extends D3CanvasBaseNodeDatum,
	LinkDatum extends D3CanvasBaseLinkDatum<NodeDatum>,
	NodeElement extends Element,
	LinkElement extends Element
>({ 
	ref, 
	nodeRenderer, 
	linkRenderer, 
	...props 
}: React.SVGProps<SVGSVGElement> & { 
	ref?: React.Ref<D3CanvasRefType<NodeDatum, LinkDatum, NodeElement, LinkElement>>, 
	nodeRenderer: D3CanvasNodeRenderer<NodeDatum, NodeElement>, 
	linkRenderer: D3CanvasLinkRenderer<LinkDatum, LinkElement> 
}) => {
	const svgElementRef = useRef<D3CanvasRefType<NodeDatum, LinkDatum, NodeElement, LinkElement>>(null);
	useIsomorphicLayoutEffect(() => {
		if(ref == null) return;
		if(typeof ref == "function") {
			ref(svgElementRef.current);
			return () => { ref(null); };
		}
		ref.current = svgElementRef.current;
		return () => { ref.current = null; };
	}, []);
	useIsomorphicLayoutEffect(() => {
		const svgElement = svgElementRef.current!;
		const handler = newD3CanvasHandler<NodeDatum, LinkDatum, NodeElement, LinkElement>(svgElement);
		svgElement.handler = handler;
		return () => { handler.container.remove(); };
	}, []);
	useIsomorphicLayoutEffect(() => {
		svgElementRef.current!.handler.nodeRenderer = nodeRenderer;
	}, [nodeRenderer]);
	useIsomorphicLayoutEffect(() => {
		svgElementRef.current!.handler.linkRenderer = linkRenderer;
	}, [linkRenderer]);
	return (
		<svg ref={svgElementRef} {...props}/>
	);
};

const splashTexts = [
	"🌍🔍 Temukan Semua 720 Elemen dari 4 Unsur Dasar!",
	"🧪💡 BFS dan DFS Siap Menguak Resep Alkimia!",
	"🌊🔥💨🌱 Dari Dasar Menuju Keajaiban—Gabungkan dan Temukan!",
	"🧠⚙️ Strategi Algoritma Bertemu Dunia Alkimia!",
	"🧬✨ Temukan Resep Tersembunyi dengan Sekali Klik!",
	"🌀🚀 Pilih DFS atau BFS—Raih Elemen Impianmu!",
	"🌳🧩 Visualisasi Pohon Resep yang Seru dan Informatif!",
	"⏱️📊 Ukur Kecepatanmu—Cek Waktu dan Node yang Dilalui!",
	"🔁🎯 Temukan Banyak Resep dalam Sekejap dengan Multithreading!",
	"🕹️🧙‍♂️ Jadi Alkemis Digital Terbaik di Dunia Little Alchemy 2!",
	"INI TUGASS BESAAR WOEEEEE"
];

export default function Page() {
	const [splashText, setSplashText] = useState<string | null>(null);
	useIsomorphicLayoutEffect(() => setSplashText(splashTexts[Math.floor(Math.random() * splashTexts.length)]), []);
	type NodeDatum = D3CanvasBaseNodeDatum & { name: string, description: string, image: string };
	type LinkDatum = D3CanvasBaseLinkDatum<NodeDatum> & {};
	type NodeElement = SVGForeignObjectElement;
	type LinkElement = SVGLineElement;
	const _rerender = useState(0)[1];
	const rerender = useCallback(() => _rerender(c => c + 1), []);
	const d3CanvasRef = useRef<D3CanvasRefType<NodeDatum, LinkDatum, NodeElement, LinkElement>>(null);
	const NodeRenderer = useCallback(({ datum }: { datum: NodeDatum }) => {
		const handle = d3CanvasRef.current!.handler;
		const edgeCount = [...handle.allLinks.values()].filter(l => l.source == datum || l.target == datum).length;
		return (
			<Popover>
				<PopoverTrigger asChild>
					<img className="w-full h-full rounded-full cursor-pointer" src={datum.image} />
				</PopoverTrigger>
				<PopoverContent>
					<div className="grid gap-4">
						<div className="space-y-2">
							<h4 className="font-medium leading-none">{datum.name}</h4>
							<p className="text-sm text-muted-foreground">{datum.description}</p>
						</div>
						<div className="grid gap-2">
							<div className="grid grid-cols-5 items-center gap-4">
								<Label htmlFor="timeTaken" className="col-span-2">Time Taken</Label>
								<Input id="timeTaken" className="col-span-3 h-8 opacity-100!" value="100ms" disabled />
							</div>
							<div className="grid grid-cols-5 items-center gap-4">
								<Label htmlFor="edgeCount" className="col-span-2">Edge Count</Label>
								<Input id="edgeCount" className="col-span-3 h-8 opacity-100!" value={edgeCount} disabled />
							</div>
							<div className="grid grid-cols-5 items-center gap-4">
								<Label htmlFor="edgeCount" className="col-span-2">All Recipes</Label>
								<Dialog>
									<DialogTrigger asChild>
										<Input id="edgeCount" className="col-span-3 h-8 opacity-100! bg-primary text-primary-foreground shadow-xs hover:bg-primary/90" type="button" value={`5 Recipes`} />
									</DialogTrigger>
									<DialogContent>
										<DialogHeader>
											<DialogTitle>{datum.name}</DialogTitle>
											<DialogDescription>{datum.description}</DialogDescription>
										</DialogHeader>
									</DialogContent>
								</Dialog>
							</div>
						</div>
					</div>
				</PopoverContent>
			</Popover>
		);
	}, []);
	const nodePortals = useMemo(() => new Map<NodeElement, [NodeDatum, React.JSX.Element]>(), []);
	const nodeRenderer = useMemo(() => {
		return {
			enter: nodes => {
				let hasChanged = false;
				const result = nodes.append("foreignObject")
					.attr("class", "rounded-full outline-none")
					.attr("width", 50)
					.attr("height", 50)
					.attr("x", -25)
					.attr("y", -25)
					.each((d, i, g) => {
						nodePortals.set(g[i], [d, <NodeRenderer datum={d} />]);
						hasChanged = true;
					});
				if(hasChanged)
					rerender();
				return result;
			},
			update: nodes => {
				nodes
					.attr("x", d => d.x! - 25)
					.attr("y", d => d.y! - 25);
			},
			exit: nodes => {
				let hasChanged = false;
				nodes.each((_, i, g) => {
					nodePortals.delete(g[i]);
					hasChanged = true;
				});
				if(hasChanged)
					rerender();
			}
		} as D3CanvasNodeRenderer<NodeDatum, NodeElement>;
	}, []);
	const linkRenderer = useMemo(() => {
		return {
			enter: links => {
				return links.append("line")
					.attr("stroke", "black")
					.attr("stroke-width", 2);
			},
			update: links => {
				links
					.attr("x1", d => d.source.x!)
					.attr("y1", d => d.source.y!)
					.attr("x2", d => d.target.x!)
					.attr("y2", d => d.target.y!);
			},
			exit: () => {}
		} as D3CanvasLinkRenderer<LinkDatum, LinkElement>;
	}, []);
	useIsomorphicLayoutEffect(() => {
		const handle = setTimeout(() => {
			const d3CanvasHandle = d3CanvasRef.current!.handler;
			d3CanvasHandle.allNodes.set("1", {
				id: "1",
				name: "Fire",
				description: "Lorem ipsum",
				image: "/Air.png"
			});
			d3CanvasHandle.allNodes.set("2", {
				id: "2",
				name: "Fire",
				description: "Lorem ipsum",
				image: "/Air.png"
			});
			d3CanvasHandle.allLinks.set("1-2", {
				id: "1-2",
				source: d3CanvasHandle.allNodes.get("1")!,
				target: d3CanvasHandle.allNodes.get("2")!,
			});
			d3CanvasHandle.refreshData();
		}, 100);
		return () => {
			clearTimeout(handle);
		};
	}, []);
	return (
		<div className="h-full flex flex-col">
			<div className="container md:h-16 p-4 flex flex-col items-start justify-between space-y-2 sm:flex-row sm:items-center sm:space-y-0">
				<h2 className="text-lg font-semibold">Nuggets</h2>
				<div className="w-full ml-auto flex space-x-2 sm:justify-end">{splashText}</div>
			</div>
			<Separator className="h-[2px]!" />
			<div className="container h-full px-4 py-6">
				<div className="h-full grid items-stretch gap-6 md:grid-cols-[1fr_200px]">
					<div className="flex flex-col space-y-4 md:order-2">
						<Button>Submit</Button>
					</div>
					<div className="md:order-1">
						<D3Canvas<NodeDatum, LinkDatum, NodeElement, LinkElement> 
							ref={d3CanvasRef}
							nodeRenderer={nodeRenderer}
							linkRenderer={linkRenderer}
							className="w-full min-h-[400px] md:min-h-[700px] border-2 border-input rounded-md flex bg-background text-base ring-offset-background" 
						/>
						{[...nodePortals].map(([n, [d, e]]) => createPortal(e, n, d.id))}
					</div>
				</div>
			</div>
		</div>
	);
}
