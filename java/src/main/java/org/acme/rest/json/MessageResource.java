package org.acme.rest.json;

import java.util.List;
import javax.inject.Inject;
import javax.persistence.EntityManager;
import javax.persistence.TypedQuery;
import javax.persistence.criteria.CriteriaBuilder;
import javax.persistence.criteria.CriteriaQuery;
import javax.persistence.criteria.Root;
import javax.transaction.Transactional;
import javax.ws.rs.Consumes;
import javax.ws.rs.DELETE;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.PUT;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;

@Path("/messages")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
public class MessageResource {

  @Inject
  EntityManager em;

  public MessageResource() {
  }

  @GET
  @Produces("application/json")
  public List<Message> list() {
    return getAllMessages();
  }

  @GET
  @Path("{id}")
  @Produces("application/json")
  public Message get(@PathParam("id") long id) {
    Message msg = em.find(Message.class, id);
    return msg;
  }

  @POST
  @Transactional
  @Produces("application/json")
  public List<Message> add(Message message) {
    em.persist(message);
    return getAllMessages();
  }

  @PUT
  @Transactional
  @Path("{id}")
  @Produces("application/json")
  public Message update(@PathParam("id") long id, Message message) {
    Message msg = em.find(Message.class, id);
    if (msg != null) {
      msg.setText(message.getText());
      em.merge(msg);
      return msg;
    }
    return null;
  }

  @DELETE
  @Transactional
  @Path("{id}")
  @Produces("application/json")
  public Message delete(@PathParam("id") long id) {
    Message msg = em.find(Message.class, id);
    if (msg != null) {
      em.remove(msg);
    }
    return msg;
  }

  private List<Message> getAllMessages() {
    CriteriaBuilder cb = em.getCriteriaBuilder();
    CriteriaQuery<Message> cq = cb.createQuery(Message.class);
    Root<Message> rootEntry = cq.from(Message.class);
    CriteriaQuery<Message> all = cq.select(rootEntry);
    TypedQuery<Message> allQuery = em.createQuery(all);
    return allQuery.getResultList();
  }
}
